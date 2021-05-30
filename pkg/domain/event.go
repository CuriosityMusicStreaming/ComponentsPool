package domain

type Event interface {
	ID() string
}

type EventHandler interface {
	Handle(event Event) error
}

type EventDispatcher interface {
	Dispatch(event Event) error
}

type EventSource interface {
	Subscribe(handler EventHandler)
}

type EventPublisher interface {
	EventDispatcher
	EventSource
}

func NewEventPublisher() EventPublisher {
	return &eventPublisher{}
}

type eventPublisher struct {
	subscribers []EventHandler
}

func (e *eventPublisher) Dispatch(event Event) error {
	for _, subscriber := range e.subscribers {
		err := subscriber.Handle(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *eventPublisher) Subscribe(handler EventHandler) {
	e.subscribers = append(e.subscribers, handler)
}
