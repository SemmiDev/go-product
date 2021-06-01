package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	AppPort, err := strconv.Atoi(os.Getenv("AppPort"))
	HttpRateLimitRequest, err := strconv.Atoi(os.Getenv("HttpRateLimitRequest"))
	HttpRateLimitTime,err := time.ParseDuration(os.Getenv("HttpRateLimitTime"))
	JwtSecretKey := os.Getenv("JwtSecretKey")
	JwtTTL, err := time.ParseDuration(os.Getenv("JwtTTL"))
	PaginationLimit, err := strconv.Atoi(os.Getenv("PaginationLimit"))
	MysqlUser := os.Getenv("MysqlUser")
	MysqlPassword := os.Getenv("MysqlPassword")
	MysqlHost := os.Getenv("MysqlHost")
	MysqlPort,err := strconv.Atoi(os.Getenv("MysqlPort"))
	MysqlDatabase := os.Getenv("MysqlDatabase")
	MysqlMaxIdleConns,err := strconv.Atoi(os.Getenv("MysqlMaxIdleConns"))
	MysqlMaxOpenConns,err := strconv.Atoi(os.Getenv("MysqlMaxOpenConns"))
	MysqlConnMaxLifetime,err := time.ParseDuration(os.Getenv("MysqlConnMaxLifetime"))
	RedisPassword := os.Getenv("RedisPassword")
	RedisHost := os.Getenv("RedisHost")
	RedisPort,err := strconv.Atoi(os.Getenv("RedisPort"))
	RedisDatabase, err := strconv.Atoi(os.Getenv("RedisDatabase"))
	RedisPoolSize, err := strconv.Atoi(os.Getenv("RedisPoolSize"))
	RedisTTL,err := time.ParseDuration(os.Getenv("RedisTTL"))

	if err != nil {
		log.Fatal("error in get value from .env")
	}

	return Config{
		AppPort:              AppPort,
		HttpRateLimitRequest: HttpRateLimitRequest,
		HttpRateLimitTime:    HttpRateLimitTime,
		JwtSecretKey:         JwtSecretKey,
		JwtTTL:               JwtTTL,
		PaginationLimit:      PaginationLimit,
		MysqlUser:            MysqlUser,
		MysqlPassword:        MysqlPassword,
		MysqlHost:            MysqlHost,
		MysqlPort:            MysqlPort,
		MysqlDatabase:        MysqlDatabase,
		MysqlMaxIdleConns:    MysqlMaxIdleConns,
		MysqlMaxOpenConns:    MysqlMaxOpenConns,
		MysqlConnMaxLifetime: MysqlConnMaxLifetime,
		RedisPassword:        RedisPassword,
		RedisHost:            RedisHost,
		RedisPort:            RedisPort,
		RedisDatabase:        RedisDatabase,
		RedisPoolSize:        RedisPoolSize,
		RedisTTL:             RedisTTL,
	}
}

var config = load()

func Cfg() *Config { return &config }