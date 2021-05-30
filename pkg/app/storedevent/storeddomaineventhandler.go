package storedevent

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/domain"
)

type EventSerializer interface {
	Serialize(event domain.Event) (string, error)
}

func NewStoredDomainEventHandler(eventStore Store, eventSerializer EventSerializer) domain.EventHandler {
	return &storedDomainEventHandler{
		eventStore:      eventStore,
		eventSerializer: eventSerializer,
	}
}

type storedDomainEventHandler struct {
	eventStore      Store
	eventSerializer EventSerializer
}

func (handler *storedDomainEventHandler) Handle(event domain.Event) error {
	body, err := handler.eventSerializer.Serialize(event)
	if err != nil {
		return err
	}

	return handler.eventStore.Append(NewStoredEvent(event.ID(), body))
}
