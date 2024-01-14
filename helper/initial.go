package helper

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/utils"
)

func InitFiles() {
	initialFile := "initialData/initialFiles.json"
	content, err := os.ReadFile(initialFile)
	if err != nil {
		log.Fatal("Error when read file: ", err)
	}

	var payload models.Files
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error when unmarshal data: ", err)
	}
	for _, file := range payload.Files {
		var content string
		if file.IsJSON {
			j, err := json.Marshal(file.Template)
			if err == nil {
				content = string(j)
			}
		} else {
			content = file.Template.(string)
		}
		err = utils.WriteToFileIfNotExists(content, "data/"+file.Name)
		if err != nil {
			log.Fatal("Error when write file: ", err)
		}
	}
}

func InitConfig() {
	configFile := "data/config.json"
	configInitialFile := "initialData/initialConfig.json"

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		if err = utils.CopyFile(configInitialFile, configFile); err != nil {
			log.Fatal("Error when copy file: ", err)
		}
	}

	configContent, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	configInitialContent, err := os.ReadFile(configInitialFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var configJson map[string]interface{}
	err = json.Unmarshal(configContent, &configJson)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var initialJson map[string]interface{}
	err = json.Unmarshal(configInitialContent, &initialJson)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	for k, v := range initialJson {
		if _, ok := configJson[k]; !ok {
			configJson[k] = v
		}
	}

	content, err := json.MarshalIndent(configJson, "", "  ")
	if err != nil {
		log.Fatal("Error when marshal data: ", err)
	}

	err = utils.WriteToFileIfNotExists(string(content), configFile)
	if err != nil {
		log.Fatal("Error when write file: ", err)
	}
}
