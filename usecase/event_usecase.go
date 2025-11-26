package usecase

import (
	"context"
	"t-meeting-backend/domain"
	"t-meeting-backend/repository"

	"github.com/google/uuid"
)

type EventUsecase interface {
	Create(ctx context.Context, e *domain.Event) error
	GetAll(ctx context.Context) ([]*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, id uuid.UUID, e *domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type eventUsecase struct {
	repo repository.EventRepository
}

func (u *eventUsecase) Create(ctx context.Context, e *domain.Event) error {
	return u.repo.Create(ctx, e)
}

func (u *eventUsecase) GetAll(ctx context.Context) ([]*domain.Event, error) {
	return u.repo.GetAll(ctx)
}
func (u *eventUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	e, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, err
	}
	return e, nil
}

func (u *eventUsecase) Update(ctx context.Context, id uuid.UUID, e *domain.Event) error {
	return u.repo.Update(ctx, id, e)
}
func (u *eventUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func NewEventUsecase(eventRepository repository.EventRepository) EventUsecase {
	return &eventUsecase{
		repo: eventRepository,
	}
}
