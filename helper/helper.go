package helper

import (
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/utils"
)

func LoadConfig(fileName string) (*models.Config, error) {
	data, err := utils.ReadJson[models.Config](fileName)
	if err != nil {
		return nil, err
	}
	return data, err
}
