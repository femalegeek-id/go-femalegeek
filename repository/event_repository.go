package repository

import (
	"context"
	"errors"
	"femalegeek/config"
	"femalegeek/repository/model"

	"github.com/go-redsync/redsync/v2"
	"github.com/jinzhu/gorm"
	"github.com/kumparan/cacher"
	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

type (
	eventRepo struct {
		db          *gorm.DB
		cacheKeeper cacher.Keeper
	}
)

// NewEventRepository create new repository
func NewEventRepository(d *gorm.DB, k cacher.Keeper) model.EventRepository {
	return &eventRepo{
		db:          d,
		cacheKeeper: k,
	}
}

// FindByID find object with specific id
func (r *eventRepo) FindByID(ctx context.Context, id int64) (event *model.Event, err error) {
	var (
		logger = log.WithFields(log.Fields{
			"context": utils.DumpIncomingContext(ctx),
			"id":      id})
		cacheKey = model.NewEventCacheKeyByID(id)
		e        model.Event
	)

	if !config.DisableCaching() {
		event, mu, err := r.findFromCacheByKey(cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		if mu == nil {
			return event, nil
		}

		defer cacher.SafeUnlock(mu)
	}

	err = r.db.Take(&e, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = r.cacheKeeper.StoreNil(cacheKey)
			if err != nil {
				logger.Error(err)
			}
			return nil, nil
		}
		logger.Error(err)
		return nil, err
	}

	event = &e

	err = r.cacheKeeper.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.ToByte(event)))
	if err != nil {
		logger.Error(err)
	}

	return event, nil
}

func (r *eventRepo) findFromCacheByKey(key string) (e *model.Event, mu *redsync.Mutex, err error) {
	reply, mu, err := r.cacheKeeper.GetOrLock(key)
	if err != nil {
		return
	}

	if reply == nil {
		return
	}

	e, err = model.NewEventFromInterface(reply)

	return
}
