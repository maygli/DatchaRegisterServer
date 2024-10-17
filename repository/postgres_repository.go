package repository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatchaSqlRepository struct {
	database   *gorm.DB
	connection *DatchaSqlConnectionCfg
	notifier   servercommon.INotifier
}

func NewSqlRepository(notifier servercommon.INotifier) (*DatchaSqlRepository, error) {
	repo := DatchaSqlRepository{
		notifier: notifier,
	}
	var err error
	repo.connection, err = NewSqlConnectionCfg()
	if err != nil {
		return nil, err
	}
	connStr := repo.connection.toString()

	repo.database, err = gorm.Open(postgres.Open(connStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	// Migrate the schema
	repo.database.AutoMigrate(&datamodel.User{})
	repo.database.AutoMigrate(&datamodel.Project{})
	repo.database.AutoMigrate(&datamodel.Device{})
	repo.database.AutoMigrate(&datamodel.Channel{})
	repo.database.AutoMigrate(&datamodel.ChannelData{})
	repo.database.AutoMigrate(&datamodel.ProjectCardItem{})
	repo.database.AutoMigrate(&datamodel.Command{})
	/*	repo.database, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}
		err = repo.database.Ping()
		if err != nil {
			return nil, err
		}
		_, err = repo.database.Exec(SQL_CREATE_USER_TABLE)
		if err != nil {
			return nil, err
		}*/
	project, _ := repo.FindProject(57)
	fmt.Println(project)
	return &repo, nil
}
