package config

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// HTTPPort get port
func HTTPPort() string {
	return viper.GetString("http_port")
}

// Env get Httpenvironment
func Env() string {
	return viper.GetString("env")
}

// LogLevel get log level
func LogLevel() string {
	return viper.GetString("log_level")
}

// DisableCaching get disable cache
func DisableCaching() bool {
	return viper.GetBool("disable_caching")
}

// CacheTTL get default ttl
func CacheTTL() time.Duration {
	if !viper.IsSet("cache_ttl") {
		return DefaultCacheTTL
	}

	return time.Duration(viper.GetInt("cache_ttl")) * time.Millisecond
}

// RedisCacheHost get redis cache host
func RedisCacheHost() string {
	return viper.GetString("redis.cache_host")
}

// RedisLockHost get redis lock host
func RedisLockHost() string {
	return viper.GetString("redis.lock_host")
}

// DatabaseDSN get  dsn
func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		PostgresUsername(),
		PostgresPassword(),
		PostgresHost(),
		PostgresDatabase())
}

// PostgresHost get postgres host
func PostgresHost() string {
	return viper.GetString("postgres.host")
}

// PostgresDatabase get postgres database
func PostgresDatabase() string {
	return viper.GetString("postgres.database")
}

// PostgresUsername get postgres username
func PostgresUsername() string {
	return viper.GetString("postgres.username")
}

// PostgresPassword get postgres password
func PostgresPassword() string {
	return viper.GetString("postgres.password")
}

// PostgresMaxIdleConns :nodoc:
func PostgresMaxIdleConns() int {
	if viper.GetInt("postgres.max_idle_conns") <= 0 {
		return DefaultPostgresMaxIdleConns
	}
	return viper.GetInt("postgres.max_idle_conns")
}

// PostgresMaxOpenConns :nodoc:
func PostgresMaxOpenConns() int {
	if viper.GetInt("postgres.max_open_conns") <= 0 {
		return DefaultPostgresMaxOpenConns
	}
	return viper.GetInt("postgres.max_open_conns")
}

// PostgresConnMaxLifetime :nodoc:
func PostgresConnMaxLifetime() time.Duration {
	if !viper.IsSet("postgres.conn_max_lifetime") {
		return DefaultPostgresConnMaxLifetime
	}
	return time.Duration(viper.GetInt("postgres.conn_max_lifetime")) * time.Millisecond
}

// PostgresPingInterval :nodoc:
func PostgresPingInterval() time.Duration {
	if viper.GetInt("postgres.ping_interval") <= 0 {
		return DefaultPostgresPingInterval
	}
	return time.Duration(viper.GetInt("postgres.ping_interval")) * time.Millisecond
}

// GetConf read configuration
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("./../../..")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("svc")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Warningf("%v", err)
	}
}
