package repository

import (
	"context"
	"errors"

	"t-meeting-backend/domain"

	"github.com/google/uuid"
)

type EventRepository interface {
	Create(c context.Context, e *domain.Event) error
	GetAll(c context.Context) ([]*domain.Event, error)
	GetByID(c context.Context, id uuid.UUID) (*domain.Event, error)
	Update(c context.Context, id uuid.UUID, e *domain.Event) error
	Delete(c context.Context, id uuid.UUID) error
}

type eventRepository struct {
	database map[uuid.UUID]*domain.Event
}

func NewEventRepository() EventRepository {
	db := make(map[uuid.UUID]*domain.Event)
	return &eventRepository{database: db}
}

func (er *eventRepository) Create(_ context.Context, e *domain.Event) error {
	id := uuid.New()
	e.ID = id
	er.database[id] = e
	return nil
}

func (er *eventRepository) GetAll(_ context.Context) ([]*domain.Event, error) {
	var res []*domain.Event
	for _, v := range er.database {
		res = append(res, v)
	}
	return res, nil
}

func (er *eventRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Event, error) {
	event := er.database[id]
	if event == nil {
		return nil, errors.New("мероприятие не найдено")
	}
	return er.database[id], nil
}

func (er *eventRepository) Update(_ context.Context, id uuid.UUID, e *domain.Event) error {
	e.ID = id
	er.database[id] = e
	return nil
}

func (er *eventRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(er.database, id)
	return nil
}
