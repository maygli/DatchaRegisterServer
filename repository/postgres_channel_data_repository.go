package repository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

const (
	LIMIT_BY_COUNT_QUERY = "DELETE FROM %s WHERE %s IN(SELECT %s FROM %s WHERE %s=%d ORDER BY %s ASC LIMIT %d);"
	ID_FIELD             = "id"
	CHANNEL_ID_FIELD     = "channel_id"
	CREATED_AT_FIELD     = "created_at"
)

func (repo *DatchaSqlRepository) AddChannelData(channelId uint, data string, unit string) error {
	chData := datamodel.ChannelData{
		Data:      data,
		Unit:      unit,
		ChannelID: channelId,
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(&chData)
	if result.Error != nil {
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	err := repo.LimitChannelByTime(channelId, 600)
	if err != nil {
		repo.database.Rollback()
	}
	return nil
}

func (repo *DatchaSqlRepository) GetAllChannelData(channelId uint) ([]datamodel.ChannelData, error) {
	var chData []datamodel.ChannelData
	result := repo.database.Where("channel_id=?", channelId).Order("created_at asc").Find(&chData)
	return chData, result.Error
}

func (repo *DatchaSqlRepository) GetLastChannelData(channelId uint) (datamodel.ChannelData, error) {
	var chData datamodel.ChannelData
	result := repo.database.Where("channel_id=?", channelId).Last(&chData)
	return chData, result.Error
}

func (repo *DatchaSqlRepository) LimitChannelByCount(channelId uint, limit int64) error {
	var currSize int64
	res := repo.database.Model(&datamodel.ChannelData{}).Where(&datamodel.ChannelData{ChannelID: channelId}).Count(&currSize)
	if res.Error != nil {
		message := fmt.Sprintf("Error get channel data size. ChannelId=%d, Error=%s", channelId, res.Error.Error())
		log.Println(message)
	}
	if currSize < limit {
		return nil
	}
	tableName := repo.getChannelDataTableName()
	execStr := fmt.Sprintf(LIMIT_BY_COUNT_QUERY, tableName, ID_FIELD, ID_FIELD, tableName, CHANNEL_ID_FIELD, channelId, CREATED_AT_FIELD, currSize-limit)
	log.Println(execStr)
	//	datas := []datamodel.ChannelData{}
	//	itemsToDel := currSize - limit
	res = repo.database.Exec(execStr)
	//	res = repo.database.Where().Delete(&datamodel.ChannelData{})
	//	res = repo.database.Model(&datamodel.ChannelData{}).Where(&datamodel.ChannelData{ChannelID: channelId}).
	//		Order("created_at ASC").Limit(int(itemsToDel)).Delete(&datamodel.ChannelData{})
	if res.Error != nil {
		message := fmt.Sprintf("Error remove channel data. ChannelId=%d, Error=%s", channelId, res.Error.Error())
		log.Println(message)
		return res.Error
	}
	return nil
}

func (repo *DatchaSqlRepository) LimitChannelByTime(channelId uint, limitSecs int64) error {
	timeLim := time.Now()
	timeLim = timeLim.Add(-time.Duration(limitSecs) * time.Second)
	//	tableName := repo.getChannelDataTableName()
	//	res := repo.database.Exec("DELETE FROM $1 WHERE created_at < $2 AND channel_id=$3", tableName, timeLim, channelId)
	res := repo.database.Where("created_at < ?", timeLim).Delete(&datamodel.ChannelData{ChannelID: channelId})
	return res.Error
}

func (repo *DatchaSqlRepository) getChannelDataTableName() string {
	stmt := &gorm.Statement{DB: repo.database}
	stmt.Parse(&datamodel.ChannelData{})
	return stmt.Schema.Table
}
