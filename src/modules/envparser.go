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
	config_path string
}

func get_envparser() *envparser {
	return &envparser{config_path: "/etc/ec2_connect_config.json"}
}

func (ep *envparser) get_credentials() (string, string, error) {
	// Reference : https://tutorialedge.net/golang/parsing-json-with-golang/
	jsonFile, err := os.Open(ep.config_path)
	if err != nil {
		log.Fatal("Unable to open /etc/ec2_connect_config.json. Abort")
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	secret_data := result["CONFIG"].(map[string]interface{})

	err = check_keys(secret_data, []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"})
	if err != nil {
		return "", "", err
	} else {
		log.Println("Found credential in configuration file")
		access_id := secret_data["AWS_ACCESS_KEY_ID"].(string)
		secret_key := secret_data["AWS_SECRET_ACCESS_KEY"].(string)
		return access_id, secret_key, nil

	}
}

func check_keys(secret_data map[string]interface{}, data []string) error {
	for val := range data {
		if _, ok := secret_data[data[val]]; !ok {
			return errors.New(fmt.Sprint("Cannot found key in configuration file : ", data[val]))
		}
	}
	return nil
}
