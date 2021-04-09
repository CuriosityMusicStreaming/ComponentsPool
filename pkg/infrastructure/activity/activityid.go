package activity

import "github.com/google/uuid"

type ID uuid.UUID

func NewActivityID() ID {
	return ID(uuid.New())
}

func ParseActivityID(value string) (ID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return [16]byte{}, err
	}
	return ID(id), nil
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}
