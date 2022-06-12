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
	userRepo struct {
		db          *gorm.DB
		cacheKeeper cacher.Keeper
	}
)

// NewUserRepository create new repository
func NewUserRepository(d *gorm.DB, k cacher.Keeper) model.UserRepository {
	return &userRepo{
		db:          d,
		cacheKeeper: k,
	}
}

// FindByID find object with specific id
func (r *userRepo) FindByID(ctx context.Context, id int64) (user *model.User, err error) {
	var (
		logger = log.WithFields(log.Fields{
			"context": utils.DumpIncomingContext(ctx),
			"id":      id})
		cacheKey = model.NewUserCacheKeyByID(id)
		u        model.User
	)

	user, mu, err := r.findFromCacheByKey(cacheKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !config.DisableCaching() && mu == nil {
		return user, nil
	}

	defer func() {
		if mu != nil {
			mu.Unlock()
		}
	}()

	err = r.db.Take(&u, "id = ?", id).Error
	if err != nil {
		log.Info("err = ", err)
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

	user = &u

	err = r.cacheKeeper.Store(mu, cacher.NewItem(cacheKey, utils.ToByte(user)))
	if err != nil {
		logger.Error(err)
	}

	return user, nil
}

func (r *userRepo) findFromCacheByKey(key string) (u *model.User, mu *redsync.Mutex, err error) {
	reply, mu, err := r.cacheKeeper.GetOrLock(key)
	if err != nil {
		return
	}

	if reply == nil {
		return
	}

	u, err = model.NewUserFromInterface(reply)

	return
}
