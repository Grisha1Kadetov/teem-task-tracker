package config

import (
	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/pkg/log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port      string `env:"PORT" env-default:"8080"`
	JWTSecret string `env:"JWT_SECRET" env-default:"secret"`
	RetryCount	int    `env:"RETRY" env-default:"3"`
	MySQL     MySQL
	Redis     Redis
}

type MySQL struct {
	Host     string `env:"MYSQL_HOST" env-default:"localhost"`
	Port     string `env:"MYSQL_PORT" env-default:"3306"`
	Database string `env:"MYSQL_DATABASE" env-default:"task_tracker"`
	User     string `env:"MYSQL_USER" env-default:"task_tracker"`
	Password string `env:"MYSQL_PASSWORD" env-default:"task_tracker_password"`
}

type Redis struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port string `env:"REDIS_PORT" env-default:"6379"`
}

func LoadConfig(l log.Logger) (*Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		l.Warn("can't load .env file")
	}
	err = cleanenv.ReadEnv(&cfg)
	return &cfg, err
}

func (cfg *Config) GetMySQLDSN() string {
	return cfg.MySQL.User + ":" + cfg.MySQL.Password + "@tcp(" + cfg.MySQL.Host + ":" + cfg.MySQL.Port + ")/" + cfg.MySQL.Database
}
