package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

type (
	// EventType :nodoc:
	EventType string

	// Event :nodoc:
	Event struct {
		ID          int64
		Slug        string
		Title       string
		Description string
		Cover       string
		Organizer   string
		StartsAt    time.Time
		EndsAt      time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   *time.Time
	}
)

// TableName :nodoc:
func (Event) TableName() string {
	return "events"
}

// EventRepository :nodoc:
type EventRepository interface {
	FindByID(ctx context.Context, id int64) (*Event, error)
}

// EventUsecase :nodoc:
type EventUsecase interface {
	FindByID(ctx context.Context, id int64) (*Event, error)
}

// NewEventCacheKeyByID :nodoc:
func NewEventCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:event:id:%d", id)
}

// NewEventFromInterface converts interface to event
func NewEventFromInterface(i interface{}) (u *Event, err error) {
	bt, _ := i.([]byte)

	err = json.Unmarshal(bt, &u)
	if err != nil {
		log.WithField("i", utils.Dump(i)).Error(err)
	}

	return
}
