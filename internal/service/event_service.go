package service

import (
	"context"
	"fmt"
	"t-meeting-backend/internal/domain"

	"github.com/google/uuid"
)

type EventService interface {
	Create(ctx context.Context, e *domain.NewEvent) (uuid.UUID, error)
	GetAll(ctx context.Context) ([]*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetPublishedByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
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
		return nil, err
	}
	return e, nil
}

func (es *eventService) GetPublishedByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := es.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get published event by ID: %w", err)
	}
	if event.Status != domain.EventStatusPublished {
		return nil, fmt.Errorf("get published event by ID: %w", domain.ErrPublishedEventNotFound)
	}

	return event, nil
}

func (es *eventService) Update(ctx context.Context, id uuid.UUID, e *domain.Event) error {
	return es.repo.Update(ctx, id, e)
}
func (es *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return es.repo.Delete(ctx, id)
}
