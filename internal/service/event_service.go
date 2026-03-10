package service

import (
	"context"
	"t-meeting-backend/internal/domain"

	"github.com/google/uuid"
)

type EventService interface {
	Create(ctx context.Context, e *domain.NewEvent) (uuid.UUID, error)
	GetAll(ctx context.Context) ([]*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, id uuid.UUID, e *domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type eventRepository interface {
	Create(ctx context.Context, e *domain.NewEvent) (uuid.UUID, error)
	GetAll(ctx context.Context) ([]*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, id uuid.UUID, e *domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type eventService struct {
	repo eventRepository
}

func NewEventService(repo eventRepository) EventService {
	return &eventService{repo: repo}
}

func (es *eventService) Create(ctx context.Context, e *domain.NewEvent) (uuid.UUID, error) {
	return es.repo.Create(ctx, e)
}

func (es *eventService) GetAll(ctx context.Context) ([]*domain.Event, error) {
	return es.repo.GetAll(ctx)
}
func (es *eventService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	e, err := es.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil
	}
	if e == nil {
		return nil, nil
	}
	return e, nil
}

func (es *eventService) Update(ctx context.Context, id uuid.UUID, e *domain.Event) error {
	return es.repo.Update(ctx, id, e)
}
func (es *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return es.repo.Delete(ctx, id)
}
