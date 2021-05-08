package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// TODO : Replace this package to viper
type YamlParserV1 struct {
	Version string `yaml:"version"`
	Spec    struct {
		Region string `yaml:"region"`
		Credentials struct {
			AccessKey string `yaml:"accessKey"`
			SecretKey string `yaml:"secretKey"`
		} `yaml:"credentials"`

		PrivateKeys []struct {
			Name string `yaml:"name"`
			Path string `yaml:"path"`
		} `yaml:"privateKeys"`

		InstanceGroups []struct {
			GroupName    string   `yaml:"name"`
			InstanceName []string `yaml:"instances"`
		} `yaml:"instanceGroups"`
	}
}

// Envparser : It parses configuration file in /etc/ec2_connect_config.json
type Envparser struct {
	configPath  string
	yamlContent YamlParserV1
}

// GetEnvparser : It returns envparser instance
func GetEnvparser() *Envparser {
	var envParser = Envparser{configPath: "/etc/ec2_connect_config.yaml"}
	envParser.OpenConfigFile()
	return &envParser
}

// GetCredentials : Return AWS credentials from file
func (ep *Envparser) GetCredentials() (string, string, error) {
	// TODO : Check yamlContent if keys exist
	accessID := ep.yamlContent.Spec.Credentials.AccessKey
	secretKey := ep.yamlContent.Spec.Credentials.SecretKey
	return accessID, secretKey, nil
}

// GetDefaultKey : Return default SSH key from configuration file
func (ep *Envparser) GetDefaultKey() (defaultKey string) {
	// Should I check whether EC2_SSH.. key exists? :D
	for _, element := range ep.yamlContent.Spec.PrivateKeys {
		if element.Name == "default" {
			return element.Path
		}
	}
	// TODO : return non-nil if default not exists
	log.Fatal("Cannot found default key (.spec.privateKeys -> 'default' key)")
	return defaultKey
}

// OpenConfigFile Return osfile pointer to parse configuration file
func (ep *Envparser) OpenConfigFile() (result map[string]interface{}) {
	// Reference : https://tutorialedge.net/golang/parsing-json-with-golang/
	yamlFile, err := ioutil.ReadFile(ep.configPath)
	if err != nil {
		log.Fatal("Unable to open /etc/ec2_connect_config.json. Abort")
	}
	err = yaml.Unmarshal(yamlFile, &ep.yamlContent)
	return result
}

// GetCustomKey : Return custom key path from configuration file
func (ep *Envparser) GetCustomKey(key string) (customKeyPath string) {
	for _, element := range ep.yamlContent.Spec.PrivateKeys {
		if element.Name == key {
			return element.Path
		}
	}
	log.Fatal(fmt.Sprintf("Cannot find key [%s] in configuration file", key))
	return customKeyPath
}

// GetGroupInstanceNames : Return instance names of named group in configuration file
func (ep *Envparser) GetGroupInstanceNames(groupName string) (instanceNames []string) {
	for _, element := range ep.yamlContent.Spec.InstanceGroups {
		if element.GroupName == groupName {
			return element.InstanceName
		}
	}
	log.Fatal(fmt.Sprintf("Cannot find group [%s] in configuration file", groupName))
	return instanceNames
}

func (ep *Envparser) GetRegion() (string) {
	region := ep.yamlContent.Spec.Region
	return region
}

//func checkKeys(secretData map[string]interface{}, data []string) error {
//	for val := range data {
//		if _, ok := secretData[data[val]]; !ok {
//			return errors.New(fmt.Sprint("Cannot found key in configuration file : ", data[val]))
//		}
//	}
//	return nil
//}
