package config

import "time"

var (
	//HTTPTimeout :nodoc:
	HTTPTimeout = 5 * time.Second
	// RetryAttempts :nodoc:
	RetryAttempts float64 = 5

	// DefaultCacheTTL 15 minute
	DefaultCacheTTL = 900000 * time.Millisecond

	// DefaultPostgresMaxIdleConns min connection pool
	DefaultPostgresMaxIdleConns = 2
	// DefaultPostgresMaxOpenConns max connection pool
	DefaultPostgresMaxOpenConns = 5
	// DefaultPostgresConnMaxLifetime :nodoc:
	DefaultPostgresConnMaxLifetime = 1 * time.Hour
	// DefaultPostgresPingInterval :nodoc:
	DefaultPostgresPingInterval = 1 * time.Second
)
