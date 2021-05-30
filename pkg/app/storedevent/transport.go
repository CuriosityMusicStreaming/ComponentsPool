package storedevent

type Transport interface {
	Name() string
	Send(eventType string, msgBody string) error
}
