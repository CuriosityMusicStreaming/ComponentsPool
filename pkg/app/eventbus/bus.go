package eventbus

import (
	"sort"
	"sync"
)

type EventID string

type Event interface {
	ID() EventID
}

type EventHandler func(event Event)

type Subscription struct {
	eventID  EventID
	id       uint64
	priority int
}

type BusSubscriber interface {
	Subscribe(eventID EventID, priority int, handler EventHandler) Subscription
	Unsubscribe(subscription Subscription)
}

type BusPublisher interface {
	Publish(event Event)
}

type Bus interface {
	BusSubscriber
	BusPublisher
}

type subscriptionInfo struct {
	id       uint64
	handler  EventHandler
	priority int
}

type subscriptionsInfoList []*subscriptionInfo

type bus struct {
	lock        sync.Mutex
	nextID      uint64
	subscribers map[EventID]subscriptionsInfoList
}

func NewBus() Bus {
	return &bus{
		subscribers: make(map[EventID]subscriptionsInfoList),
	}
}

func (b *bus) Subscribe(eventID EventID, priority int, handler EventHandler) Subscription {
	b.lock.Lock()
	defer b.lock.Unlock()

	id := b.nextID
	b.nextID++

	b.subscribers[eventID] = append(b.subscribers[eventID], &subscriptionInfo{
		id:       id,
		handler:  handler,
		priority: priority,
	})

	return Subscription{
		eventID:  eventID,
		id:       id,
		priority: priority,
	}
}

func (b *bus) Unsubscribe(subscription Subscription) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if subscribers, ok := b.subscribers[subscription.eventID]; ok {
		for id, info := range subscribers {
			if info.id == subscription.id {
				subscribers = append(subscribers[:id], subscribers[id+1:]...)
				break
			}
		}
		if len(subscribers) == 0 {
			delete(b.subscribers, subscription.eventID)
		} else {
			b.subscribers[subscription.eventID] = subscribers
		}
	}
}

func (b *bus) Publish(event Event) {
	infos := b.copySubscriptions(event.ID())

	sort.SliceStable(infos, func(i, j int) bool {
		return infos[i].priority > infos[j].priority
	})

	for _, sub := range infos {
		sub.handler(event)
	}
}

func (b *bus) copySubscriptions(eventID EventID) subscriptionsInfoList {
	b.lock.Lock()
	defer b.lock.Unlock()

	if infos, ok := b.subscribers[eventID]; ok {
		return infos
	}

	return subscriptionsInfoList{}
}
