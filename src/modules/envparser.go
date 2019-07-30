package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type envparser struct {
	configPath string
}

func getEnvparser() *envparser {
	return &envparser{configPath: "/etc/ec2_connect_config.json"}
}

func (ep *envparser) getCredentials() (string, string, error) {
	// Reference : https://tutorialedge.net/golang/parsing-json-with-golang/
	jsonFile, err := os.Open(ep.configPath)
	if err != nil {
		log.Fatal("Unable to open /etc/ec2_connect_config.json. Abort")
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	secretData := result["CONFIG"].(map[string]interface{})

	err = checkKeys(secretData, []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"})
	if err != nil {
		return "", "", err
	}

	accessID := secretData["AWS_ACCESS_KEY_ID"].(string)
	secretKey := secretData["AWS_SECRET_ACCESS_KEY"].(string)
	return accessID, secretKey, nil
}

func checkKeys(secretData map[string]interface{}, data []string) error {
	for val := range data {
		if _, ok := secretData[data[val]]; !ok {
			return errors.New(fmt.Sprint("Cannot found key in configuration file : ", data[val]))
		}
	}
	return nil
}
