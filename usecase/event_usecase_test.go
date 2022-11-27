package usecase

import (
	"context"
	"errors"
	"femalegeek/repository/model"
	"femalegeek/repository/model/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventUsecase_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := &model.Event{
		ID:    1,
		Title: "Femalegeek",
	}

	t.Run("Success", func(t *testing.T) {
		eventRepo := mock.NewMockEventRepository(ctrl)
		eventRepo.EXPECT().FindByID(gomock.Any(), event.ID).Times(1).Return(event, nil)

		eventUsecase := NewEventUsecase(eventRepo)

		res, err := eventUsecase.FindByID(context.TODO(), event.ID)
		assert.NoError(t, err)
		assert.EqualValues(t, event, res)
	})

	t.Run("Error from repo", func(t *testing.T) {
		eventRepo := mock.NewMockEventRepository(ctrl)
		eventRepo.EXPECT().FindByID(gomock.Any(), event.ID).Times(1).Return(nil, errors.New("error"))

		eventUsecase := NewEventUsecase(eventRepo)

		res, err := eventUsecase.FindByID(context.TODO(), event.ID)
		assert.Error(t, err)
		assert.Equal(t, errors.New("error"), err)
		assert.Nil(t, res)
	})

	t.Run("Not found", func(t *testing.T) {
		eventRepo := mock.NewMockEventRepository(ctrl)
		eventRepo.EXPECT().FindByID(gomock.Any(), event.ID).Times(1).Return(nil, nil)

		eventUsecase := NewEventUsecase(eventRepo)

		res, err := eventUsecase.FindByID(context.TODO(), event.ID)
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, res)
	})

}
