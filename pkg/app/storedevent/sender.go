package storedevent

import (
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"time"
)

type ErrorHandler func(err error)

type Sender interface {
	Increment()
	Stop()
}

type EventsDispatchTracker interface {
	TrackLastID(transportName string, id ID) error
	LastId(transportName string) (*ID, error)
	Lock() error
	Unlock() error
}

func NewStoredEventSender(eventStore Store, tracker EventsDispatchTracker, transports Transport, delay time.Duration, handler ErrorHandler) Sender {
	stopChan := make(chan struct{})
	s := &storedEventSender{
		eventStore:   eventStore,
		tracker:      tracker,
		transport:    transports,
		errorHandler: handler,
		stopChan:     stopChan,
	}
	s.wg.Add(1)
	s.start(delay)

	return s
}

type storedEventSender struct {
	wg               sync.WaitGroup
	eventStore       Store
	tracker          EventsDispatchTracker
	transport        Transport
	errorHandler     ErrorHandler
	stopChan         chan struct{}
	dispatchRequests int32
}

func (sender *storedEventSender) Increment() {
	done := false
	for !done {
		dispatchRequests := sender.dispatchRequests
		done = atomic.CompareAndSwapInt32(&sender.dispatchRequests, dispatchRequests, dispatchRequests+1)
	}
}

func (sender *storedEventSender) Stop() {
	sender.stopChan <- struct{}{}
	sender.wg.Wait()
}

func (sender *storedEventSender) start(delay time.Duration) {
	ticker := time.NewTicker(delay)

	go func() {
		for {
			select {
			case <-ticker.C:
				dispatchRequests := atomic.LoadInt32(&sender.dispatchRequests)
				if dispatchRequests > 0 {
					err := sender.dispatchEvents(dispatchRequests)
					if err != nil {
						sender.errorHandler(err)
					}
				}
			case <-sender.stopChan:
				sender.wg.Done()
				return
			}
		}
	}()
}

func (sender *storedEventSender) dispatchEvents(dispatchRequests int32) (err error) {
	err = sender.tracker.Lock()
	if err != nil {
		return err
	}

	defer func() {
		unlockErr := sender.tracker.Unlock()
		if unlockErr != nil {
			if err != nil {
				err = errors.Wrap(err, unlockErr.Error())
			} else {
				err = unlockErr
			}
		}
	}()

	lastID, err := sender.tracker.LastId(sender.transport.Name())
	if err != nil {
		return err
	}

	events, err := sender.eventStore.GetAllAfter(lastID)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		err2 := sender.transport.Send(event.Type, event.Body)
		if err2 != nil {
			return err2
		}
	}

	lastID = &events[len(events)-1].ID

	err = sender.tracker.TrackLastID(sender.transport.Name(), *lastID)
	if err != nil {
		return err
	}

	atomic.CompareAndSwapInt32(&sender.dispatchRequests, dispatchRequests, 0)

	return err
}
