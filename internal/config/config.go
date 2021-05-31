package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort int

	HttpRateLimitRequest int
	HttpRateLimitTime    time.Duration

	JwtSecretKey string
	JwtTTL       time.Duration

	PaginationLimit int

	MysqlUser            string
	MysqlPassword        string
	MysqlHost            string
	MysqlPort            int
	MysqlDatabase        string
	MysqlMaxIdleConns    int
	MysqlMaxOpenConns    int
	MysqlConnMaxLifetime time.Duration

	RedisPassword string
	RedisHost     string
	RedisPort     int
	RedisDatabase int
	RedisPoolSize int
	RedisTTL      time.Duration
}

func load() Config {
	fang := viper.New()

	fang.SetConfigFile(".env")
	fang.AddConfigPath(".")
	configLocation, available := os.LookupEnv("CONFIG_LOCATION")
	if available {
		fang.AddConfigPath(configLocation)
	}

	fang.AutomaticEnv()
	fang.ReadInConfig()

	return Config{
		AppPort:              9090,
		HttpRateLimitRequest: 100,
		HttpRateLimitTime:    time.Second,
		JwtSecretKey:         "secret",
		JwtTTL:               time.Hour * 48,
		PaginationLimit:      100,
		MysqlUser:            "sammidev",
		MysqlPassword:        "sammidev",
		MysqlHost:            "localhost",
		MysqlPort:            3306,
		MysqlDatabase:        "test",
		MysqlMaxIdleConns:    5,
		MysqlMaxOpenConns:    10,
		MysqlConnMaxLifetime: time.Minute * 30,
		RedisPassword:        "",
		RedisHost:            "localhost",
		RedisPort:            6379,
		RedisDatabase:        0,
		RedisPoolSize:        10,
		RedisTTL:             time.Hour*1,
	}
}

var config = load()

func Cfg() *Config { return &config }