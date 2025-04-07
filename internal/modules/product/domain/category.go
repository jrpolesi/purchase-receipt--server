package domain

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID
	Name string
}

func NewCategory(id uuid.UUID, name string) *Category {
	if id == uuid.Nil {
		id = uuid.New()
	}

	return &Category{
		ID:   id,
		Name: name,
	}
}

