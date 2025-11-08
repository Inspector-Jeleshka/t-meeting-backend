package usecase

import "t-meeting-backend/repository"

type EventUsecase interface {
	repository.EventRepository
}

func NewEventUsecase(eventRepository repository.EventRepository) EventUsecase {
	return eventRepository
}
