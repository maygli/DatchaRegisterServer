package repository

import (
	"datcha/servercommon"
	"errors"
	"fmt"
)

type DatchaSqlConnectionCfg struct {
	UserName string `json:"username" env:"DATCHA_SQL_USER_NAME, overwrite"`
	Password string `json:"password" env:"DATCHA_SQL_PASSWORD, overwrite"`
	Host     string `json:"host" env:"DATCHA_SQL_HOST, overwrite"`
	Port     int    `json:"port" env:"DATCHA_SQL_PORT, overwrite"`
	DbName   string `json:"database" env:"DATCHA_SQL_DATABASE, overwrite"`
}

const (
	SQL_CONFIG_FILE_NAME string = "sqlconfig.json"
)

func NewSqlConnectionCfg() (*DatchaSqlConnectionCfg, error) {
	cfg := DatchaSqlConnectionCfg{}
	err := servercommon.ConfigureServer(&cfg, SQL_CONFIG_FILE_NAME)
	if err != nil {
		return nil, errors.New("can't configure SQL connection")
	}
	err = cfg.validateConnection()
	return &cfg, err
}

func (conn DatchaSqlConnectionCfg) toString() string {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conn.Host, conn.Port, conn.UserName, conn.Password, conn.DbName)
	return connStr
}

func (conn DatchaSqlConnectionCfg) validateConnection() error {
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
