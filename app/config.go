package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	Name       string
	Dialect    string
	Datasource string
	LogStorage string
	SwagHost   string
	HostPort   int
}

func ReadConfig() *Config {
	dbPort, err := strconv.Atoi(os.Getenv("DBPORT"))
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	hostPort, err := strconv.Atoi(os.Getenv("HOSTPORT"))
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	result := Config{
		Host:       os.Getenv("DBHOST"),
		Port:       dbPort,
		User:       os.Getenv("DBUSER"),
		Password:   os.Getenv("DBPWD"),
		Name:       os.Getenv("DBNAME"),
		Dialect:    os.Getenv("DBDIALECT"),
		LogStorage: os.Getenv("LOGSTORAGE"),
		SwagHost:   os.Getenv("SWAGHOST"),
		HostPort:   hostPort,
	}

	result.Datasource = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		result.Host,
		result.Port,
		result.User,
		result.Password,
		result.Name,
	)

	return &result
}
