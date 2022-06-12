package usecase

import (
	"context"
	"femalegeek/repository/model"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo model.UserRepository
}

// NewUserUsecase :nodoc:
func NewUserUsecase(userRepo model.UserRepository) model.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// FindByID :nodoc:
func (u *userUsecase) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"ctx": utils.DumpIncomingContext(ctx),
			"id":  id,
		}).Error(err)
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	return user, nil
}
