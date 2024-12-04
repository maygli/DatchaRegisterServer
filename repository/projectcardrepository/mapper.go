package projectcardrepository

import (
	"datcha/datamodel"

	"github.com/jinzhu/copier"
)

func getNewRepProjectCardItem(item *datamodel.ProjectCardItem) (*RepProjectCardItem, error) {
	repItem := RepProjectCardItem{}
	err := copier.Copy(&repItem, item)
	return &repItem, err
}

func getNewDMProjectCardItem(item *RepProjectCardItem) (*datamodel.ProjectCardItem, error) {
	dmItem := datamodel.ProjectCardItem{}
	err := copier.Copy(&dmItem, item)
	return &dmItem, err
}

func getNewDMProjectCardItems(items []RepProjectCardItem) ([]*datamodel.ProjectCardItem, error) {
	result := make([]*datamodel.ProjectCardItem, len(items))
	for indx, item := range items {
		dmItem, err := getNewDMProjectCardItem(&item)
		if err != nil {
			return result, err
		}
		result[indx] = dmItem
	}
	return result, nil
}
