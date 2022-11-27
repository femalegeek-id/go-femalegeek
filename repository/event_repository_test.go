package repository

import (
	"context"
	"femalegeek/config"
	"femalegeek/repository/model"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/jinzhu/gorm"
	"github.com/kumparan/cacher"
	"github.com/kumparan/go-connect"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func initializeTest() {
	config.GetConf()
	setupLogger()
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		},
		Line: true,
		File: true,
	}

	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	verbose, _ := strconv.ParseBool(os.Getenv("VERBOSE"))
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func initializePostgresMockConn() (db *gorm.DB, mock sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	db, err = gorm.Open("postgres", mockDB)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return
}

func initializeEventRepository(dbGorm *gorm.DB, t *testing.T, mr *miniredis.Miniredis) (ar model.EventRepository) {
	initializeTest()

	k := cacher.NewKeeper()
	r, _ := connect.NewRedigoRedisConnectionPool("redis://"+mr.Addr(), nil)
	k.SetConnectionPool(r)
	k.SetLockConnectionPool(r)
	k.SetWaitTime(1 * time.Second) // override wait time to 1 second

	ar = NewEventRepository(dbGorm, k)

	return
}

func TestEventRepo_FindByID(t *testing.T) {
	viper.Set("disable_caching", false)
	m, _ := miniredis.Run()
	defer m.Close()

	article := &model.Event{
		ID:    123,
		Title: "Femalegeek",
	}

	t.Run("Success from db", func(t *testing.T) {
		db, dbmock := initializePostgresMockConn()
		ar := initializeEventRepository(db, t, m)
		ctx := context.TODO()

		queryResult := sqlmock.NewRows([]string{"id", "title"}).
			AddRow(article.ID, article.Title)
		dbmock.ExpectQuery("^SELECT .+ FROM \"events\"").WillReturnRows(queryResult)

		res, err := ar.FindByID(ctx, article.ID)
		assert.NoError(t, err)
		assert.Equal(t, res.ID, article.ID)
		assert.Equal(t, res.Title, article.Title)
	})

	t.Run("Success from redis", func(t *testing.T) {
		ar := initializeEventRepository(nil, t, m)
		ctx := context.TODO()

		res, err := ar.FindByID(ctx, article.ID)
		assert.NoError(t, err)
		assert.Equal(t, res.ID, article.ID)
		assert.Equal(t, res.Title, article.Title)
	})

	t.Run("Not found", func(t *testing.T) {
		db, dbmock := initializePostgresMockConn()
		ar := initializeEventRepository(db, t, m)
		ctx := context.TODO()

		dbmock.ExpectQuery("^SELECT .+ FROM \"events\"").WillReturnError(gorm.ErrRecordNotFound)

		res, err := ar.FindByID(ctx, 23456)
		assert.Nil(t, err)
		assert.Nil(t, res)
	})

	t.Run("Error from db", func(t *testing.T) {
		db, dbmock := initializePostgresMockConn()
		ar := initializeEventRepository(db, t, m)
		ctx := context.TODO()

		dbmock.ExpectQuery("^SELECT .+ FROM \"events\"").WillReturnError(gorm.ErrUnaddressable)

		res, err := ar.FindByID(ctx, 34567)
		assert.Equal(t, gorm.ErrUnaddressable, err)
		assert.Nil(t, res)
	})

}
