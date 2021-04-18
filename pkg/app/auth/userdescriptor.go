package auth

import (
	"encoding/json"
	"github.com/google/uuid"
)

type UserDescriptor struct {
	UserID uuid.UUID
}

type UserDescriptorSerializer interface {
	Serialize(UserDescriptor) (string, error)
	Deserialize(value string) (UserDescriptor, error)
}

func NewUserDescriptorSerializer() UserDescriptorSerializer {
	return &userDescriptorSerializer{}
}

type userDescriptorSerializer struct {
}

func (serializer *userDescriptorSerializer) Serialize(descriptor UserDescriptor) (string, error) {
	jsonDesc := jsonUserDescriptor{
		UserID: descriptor.UserID,
	}
	bytes, err := json.Marshal(jsonDesc)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (serializer *userDescriptorSerializer) Deserialize(value string) (UserDescriptor, error) {
	jsonDesc := jsonUserDescriptor{}
	err := json.Unmarshal([]byte(value), &jsonDesc)
	if err != nil {
		return UserDescriptor{}, err
	}

	return UserDescriptor{
		UserID: jsonDesc.UserID,
	}, err
}

type jsonUserDescriptor struct {
	UserID uuid.UUID `json:"user_id"`
}
