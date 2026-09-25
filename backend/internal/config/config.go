package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config 应用配置，全部通过环境变量注入。
type Config struct {
	AppName    string `env:"APP_NAME" envDefault:"orienteering"`
	AppEnv     string `env:"APP_ENV" envDefault:"production"`
	Port       string `env:"BACKEND_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"44004"`
	DBName     string `env:"DB_NAME" envDefault:"orienteering_db"`
	DBUser     string `env:"DB_USER" envDefault:"orienteering_user"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"orienteering_pwd"`
	JWTSecret  string `env:"JWT_SECRET" envDefault:"change_me_to_a_long_random_string"`
	JWTExpire  int    `env:"JWT_EXPIRE_HOURS" envDefault:"72"`
	RedisAddr  string `env:"REDIS_ADDR" envDefault:"localhost:46304"`
	RedisPass  string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB    int    `env:"REDIS_DB" envDefault:"0"`
	// CORS 允许来源，逗号分隔
	CORSOrigins string `env:"CORS_ORIGINS" envDefault:"*"`
}

// Load 解析环境变量为配置。
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// DSN 返回 PostgreSQL 连接串。
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}
