package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Envparser : It parses configuration file in /etc/ec2_connect_config.json
type Envparser struct {
	configPath string
}

// GetEnvparser : It returns envparser instance
func GetEnvparser() *Envparser {
	return &Envparser{configPath: "/etc/ec2_connect_config.json"}
}

// GetCredentials : Return AWS credentials from file
func (ep *Envparser) GetCredentials() (string, string, error) {
	secretData := ep.OpenConfigFile()["CONFIG"].(map[string]interface{})
	err := checkKeys(secretData, []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"})
	if err != nil {
		return "", "", err
	}

	accessID := secretData["AWS_ACCESS_KEY_ID"].(string)
	secretKey := secretData["AWS_SECRET_ACCESS_KEY"].(string)
	return accessID, secretKey, nil
}

// GetDefaultKey : Return default SSH key from configuration file
func (ep *Envparser) GetDefaultKey() (defaultKey string) {
	configData := ep.OpenConfigFile()["CONFIG"].(map[string]interface{})
	// Should I check whether EC2_SSH.. key exists? :D
	defaultKey = configData["EC2_SSH_PRIVATE_KEY_DEFAULT"].(string)
	return defaultKey
}

// OpenConfigFile Return osfile pointer to parse configuration file
func (ep *Envparser) OpenConfigFile() (result map[string]interface{}) {
	// Reference : https://tutorialedge.net/golang/parsing-json-with-golang/
	jsonFile, err := os.Open(ep.configPath)
	if err != nil {
		log.Fatal("Unable to open /etc/ec2_connect_config.json. Abort")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	jsonFile.Close()
	return result
}

// GetCustomKey : Return custom key path from configuration file
func (ep *Envparser) GetCustomKey(key string) (customKeyPath string) {
	configData := ep.OpenConfigFile()["CONFIG"].(map[string]interface{})
	if configData[key] != nil {
		customKeyPath = configData[key].(string)
	} else {
		log.Fatal(fmt.Sprintf("Cannot find key [%s] in configuration file", key))
	}
	return customKeyPath
}

// GetGroupInstanceNames : Return instance names of named group in configuration file
func (ep *Envparser) GetGroupInstanceNames(groupName string) (instanceNames []string) {
	configData := ep.OpenConfigFile()[groupName]
	if configData != nil {
		instanceGroupData := configData.([]interface{})
		instanceNames = make([]string, len(instanceGroupData))
		for index, value := range instanceGroupData {
			instanceNames[index] = value.(string)
		}
	} else {
		log.Fatal(fmt.Sprintf("Cannot find group [%s] in configuration file", groupName))
	}
	return instanceNames
}

func checkKeys(secretData map[string]interface{}, data []string) error {
	for val := range data {
		if _, ok := secretData[data[val]]; !ok {
			return errors.New(fmt.Sprint("Cannot found key in configuration file : ", data[val]))
		}
	}
	return nil
}
