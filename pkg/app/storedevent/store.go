package storedevent

import "github.com/google/uuid"

type ID uuid.UUID

type StoredEvent struct {
	ID   ID
	Type string
	Body string
}

func NewStoredEvent(Type string, body string) StoredEvent {
	return StoredEvent{
		ID:   ID(uuid.New()),
		Type: Type,
		Body: body,
	}
}

type Store interface {
	Append(event StoredEvent) error
	GetAllAfter(id *ID) ([]StoredEvent, error)
}
