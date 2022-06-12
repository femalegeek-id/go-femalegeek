package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

const (
	UserStatusPending  = UserStatus("PENDING")
	UserStatusActive   = UserStatus("ACTIVE")
	UserStatusInactive = UserStatus("INACTIVE")
	UserStatusBanned   = UserStatus("BANNED")
)

type (
	UserStatus string
	User       struct {
		ID        int64
		FullName  string
		NickName  string
		Job       string
		Company   string
		Phone     string
		Email     string
		Status    UserStatus
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
)

func (User) TableName() string {
	return "users"
}

// UserRepository :nodoc:
type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
}

// UserUsecase :nodoc:
type UserUsecase interface {
	FindByID(ctx context.Context, id int64) (*User, error)
}

// NewUserCacheKeyByID :nodoc:
func NewUserCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:user:id:%d", id)
}

// NewUserFromInterface converts interface to user
func NewUserFromInterface(i interface{}) (u *User, err error) {
	bt, _ := i.([]byte)

	err = json.Unmarshal(bt, &u)
	if err != nil {
		log.WithField("i", utils.Dump(i)).Error(err)
	}

	return
}

// ErrorResponse represent error response from kitabisa
type ErrorResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message,omitempty"`
}
