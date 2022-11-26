package db

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"femalegeek/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // GORM pq
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

var (
	// DB represents gorm DB
	DB *gorm.DB
	// StopTickerCh signal for closing ticker channel
	StopTickerCh chan bool

	sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)
)

// GormLogger :nodoc:
type GormLogger struct{}

// InitializePostgresConn :nodoc:
func InitializePostgresConn() {
	conn, err := openPostgresConn(config.DatabaseDSN())
	if err != nil {
		log.WithField("databaseDSN", config.DatabaseDSN()).Fatal("failed to connect postgres database: ", err)
	}

	DB = conn
	StopTickerCh = make(chan bool)

	go checkConnection(time.NewTicker(config.PostgresPingInterval()))

	if config.Env() == "local" || config.Env() == "test" {
		DB.LogMode(true)
	}

	log.Info("Connection to Postgres Server success...")
}

func checkConnection(ticker *time.Ticker) {
	for {
		select {
		case <-StopTickerCh:
			ticker.Stop()
			return
		case <-ticker.C:
			if err := DB.DB().Ping(); err != nil {
				reconnectPostgresConn()
			}
		}
	}
}

func reconnectPostgresConn() {
	b := backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    100 * time.Millisecond,
		Max:    1 * time.Second,
	}

	for b.Attempt() < config.RetryAttempts {
		conn, err := openPostgresConn(config.DatabaseDSN())
		if err != nil {
			log.WithField("databaseDSN", config.DatabaseDSN()).Error("failed to connect postgres database: ", err)
		}

		if conn != nil {
			DB = conn
			break
		}
		time.Sleep(b.Duration())
	}

	if b.Attempt() >= config.RetryAttempts {
		log.Fatal("maximum retry to connect database")
	}
	b.Reset()
}

func openPostgresConn(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetLogger(new(GormLogger))
	conn.DB().SetMaxIdleConns(config.PostgresMaxIdleConns())
	conn.DB().SetMaxOpenConns(config.PostgresMaxOpenConns())
	conn.DB().SetConnMaxLifetime(config.PostgresConnMaxLifetime())

	return conn, nil
}

// Print :nodoc:
func (g *GormLogger) Print(values ...interface{}) {
	switch {
	case values[0] == "sql":
		var formattedValues []interface{}
		for _, value := range values[4].([]interface{}) {
			formattedValues = append(formattedValues, g.formatValueByType(value))
		}
		log.WithFields(log.Fields{"took": values[2], "type": "sql"}).Info(fmt.Sprintf(sqlRegexp.ReplaceAllString(values[3].(string), "%v"), formattedValues...))
	case values[0] == "log":
		log.WithFields(log.Fields{"type": "gorm-log"}).Info(values[2])
	default:
		// do nothing and goodbye
	}
}

func (g *GormLogger) formatValueByType(value interface{}) string {
	indirectValue := reflect.Indirect(reflect.ValueOf(value))
	if indirectValue.IsValid() {
		value = indirectValue.Interface()
		if t, ok := value.(time.Time); ok {
			return fmt.Sprintf("'%v'", t.Format(time.RFC3339Nano))
		} else if b, ok := value.([]byte); ok {
			return fmt.Sprintf("'%v'", string(b))
		} else if r, ok := value.(driver.Valuer); ok {
			if value, err := r.Value(); err == nil && value != nil {
				return fmt.Sprintf("'%v'", value)
			}
			return "NULL"
		}
		return fmt.Sprintf("'%v'", value)
	}
	return fmt.Sprintf("'%v'", value)
}
