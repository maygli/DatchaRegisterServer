package repositorycommon

import (
	"datcha/servercommon"
	"errors"
	"fmt"
)

type SqlConnectionConfig struct {
	UserName string `json:"username" env:"${SERVER_NAME}_DB_USER_NAME"`
	Password string `json:"password" env:"${SERVER_NAME}_DB_PASSOWRD"`
	Host     string `json:"host" env:"${SERVER_NAME}_DB_HOST"`
	Port     int    `json:"port" env:"${SERVER_NAME}_DB_PORT"`
	DbName   string `json:"database" env:"${SERVER_NAME}_DB_DATABASE"`
}

const (
	SQL_CONFIG_PATH string = "$.database"
)

func NewSqlConnectionConfig(cfgReader *servercommon.ConfigurationReader) (*SqlConnectionConfig, error) {
	connection := SqlConnectionConfig{}
	err := cfgReader.ReadConfiguration(&connection, SQL_CONFIG_PATH)
	return &connection, err
}

func (conn SqlConnectionConfig) toString() string {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conn.Host, conn.Port, conn.UserName, conn.Password, conn.DbName)
	return connStr
}

func (conn SqlConnectionConfig) validateConnection() error {
	if conn.UserName == "" {
		return errors.New("SQL username should not be empty. Please fill sqlconfig.json or set environment variable DATCHA_SQL_USER_NAME")
	}
	if conn.Password == "" {
		return errors.New("SQL password should not be empty. Please fill sqlconfig.json or set environment variable DATCHA_SQL_PASSWORD")
	}
	if conn.Host == "" {
		return errors.New("Postgresql host should not be empty. Please fill sqlconfig.json or set environment variable DATCHA_SQL_HOST")
	}
	if conn.Port == 0 {
		return errors.New("Postgresql port should not be empty. Please fill sqlconfig.json or set environment variable DATCHA_SQL_PORT")
	}
	return nil
}
