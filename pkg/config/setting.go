package config

import (
	"time"
)

type StoreSetting struct {
	DSN             string        `json:"dsn" yaml:"dsn" env:"KIR_DB_DSN"`
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifeTime time.Duration `json:"conn_max_life_time" yaml:"conn_max_life_time"`
	LogLevel        string        `json:"log_level" yaml:"log_level"`
	ShowSQL         bool          `json:"show_sql" yaml:"show_sql"`
}

type GlobalSetting struct {
	Store StoreSetting `json:"store" yaml:"store"`
}
