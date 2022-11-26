package usecase

import (
	"context"
	"femalegeek/repository/model"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

type eventUsecase struct {
	eventRepo model.EventRepository
}

// NewEventUsecase :nodoc:
func NewEventUsecase(eventRepo model.EventRepository) model.EventUsecase {
	return &eventUsecase{
		eventRepo: eventRepo,
	}
}

// FindByID :nodoc:
func (e *eventUsecase) FindByID(ctx context.Context, id int64) (*model.Event, error) {
	event, err := e.eventRepo.FindByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"ctx": utils.DumpIncomingContext(ctx),
			"id":  id,
		}).Error(err)
		return nil, err
	}
	if event == nil {
		return nil, ErrNotFound
	}

	return event, nil
}
