package server

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	hubIsCreated int32 = iota
	hubIsRunning
	hubIsStopped
)

var ErrStopped = errors.New("hub is stopped by signal without errors")

type Hub interface {
	AddServer(server Server)
	Run() error
}

type hub struct {
	state           int32
	wg              sync.WaitGroup
	stopChan        chan struct{}
	stoppers        []StopFunc
	startServerChan chan Server
	errors          chan error
	reportErrorOnce sync.Once
}

func NewHub(stopChan chan struct{}) Hub {
	startServerChan := make(chan Server)

	h := &hub{
		startServerChan: startServerChan,
		stopChan:        stopChan,
		state:           hubIsCreated,
	}

	go func() {
		for server := range startServerChan {
			h.startServer(server)
		}
	}()

	return h
}

func (h *hub) AddServer(server Server) {
	h.startServerChan <- server
}

func (h *hub) Run() error {
	var err error

	// Wait for error or stopChan message and stop all servers
	h.wg.Add(1)
	go func() {
		select {
		case err = <-h.errors:
			_ = h.stop()
		case <-h.stopChan:
			err = h.stop()
			if err == nil {
				err = ErrStopped
			}
		}
		h.wg.Done()
	}()

	// Wait until all goroutines finished
	h.wg.Wait()
	return err
}

func (h *hub) startServer(server Server) {
	started := atomic.CompareAndSwapInt32(&h.state, hubIsCreated, hubIsRunning)

	h.wg.Add(1)
	if !started && atomic.LoadInt32(&h.state) == hubIsStopped {
		h.wg.Done()
		return
	}

	h.stoppers = append(h.stoppers, server.Stop)

	go func() {
		err := server.Serve()
		h.reportError(err)
		h.wg.Done()
	}()
}

func (h *hub) stop() error {
	stopped := atomic.CompareAndSwapInt32(&h.state, hubIsCreated, hubIsStopped) ||
		atomic.CompareAndSwapInt32(&h.state, hubIsRunning, hubIsStopped)
	if !stopped {
		return nil
	}

	var err error

	for _, stopper := range h.stoppers {
		stopErr := stopper()
		if err == nil && stopErr != nil {
			err = stopErr
		}
	}

	return err
}

func (h *hub) reportError(err error) {
	if atomic.LoadInt32(&h.state) == hubIsStopped {
		return
	}

	h.reportErrorOnce.Do(func() {
		h.errors <- err
	})
}
