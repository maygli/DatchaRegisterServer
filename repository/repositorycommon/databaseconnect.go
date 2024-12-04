package repositorycommon

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"datcha/servercommon"
)

type DatchaSqlRepository struct {
	database   *gorm.DB
	connection *SqlConnectionConfig
	//	notifier   servercommon.INotifier
}

func NewDatabase(cfgReader *servercommon.ConfigurationReader) (*gorm.DB, error) {
	var err error
	conn, err := NewSqlConnectionConfig(cfgReader)
	if err != nil {
		return nil, err
	}
	connStr := conn.toString()
	return gorm.Open(postgres.Open(connStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
}
