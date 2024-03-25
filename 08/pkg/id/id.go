package id

import (
	"fmt"

	"github.com/google/uuid"
)

type User uuid.UUID

func (u *User) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*u = User(id)
	return nil
}

func (u User) String() string {
	return uuid.UUID(u).String()
}

func (u *User) Scan(data any) error {
	return scanUUID((*uuid.UUID)(u), "User", data)
}

func (u User) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(u).String()), nil
}

func (u *User) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(u), "User", data)
}

func scanUUID(u *uuid.UUID, idTypeName string, data any) error {
	if err := u.Scan(data); err != nil {
		return fmt.Errorf("scanning %q id value: %w", idTypeName, err)
	}
	return nil
}

func unmarshalUUID(u *uuid.UUID, idTypeName string, data []byte) error {
	if err := u.UnmarshalText(data); err != nil {
		return fmt.Errorf("parsing %q id value: %w", idTypeName, err)
	}
	return nil
}
