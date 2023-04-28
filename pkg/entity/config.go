package entity

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type statementGlobal struct {
	ApiLayerWhoisUrl string `yaml:"apiLayerWhoisUrl"`
	Ip2WhoisUrl      string `yaml:"ip2WhoisUrl"`
	ApiLayerApiKey   string `yaml:"apiLayerApiKey"`
	WhoisXmlApiKey   string `yaml:"whoisXmlApiKey"`
	Ip2WhoisApiKey   string `yaml:"ip2WhoisApiKey"`
	NotifyMethod     string `yaml:"notifyMethod"`
	AfterMonths      int    `yaml:"afterMonths"`
}

type statementDingDing struct {
	AlertNotifyUrl  string   `yaml:"alertNotifyUrl"`
	MsgType         string   `yaml:"msgType"`
	AlertTemplate   string   `yaml:"alertTemplate"`
	AlertTitle      string   `yaml:"alertTitle"`
	UnknownTemplate string   `yaml:"unknownTemplate"`
	UnknownTitle    string   `yaml:"unknownTitle"`
	AtWho           []string `yaml:"atWho"`
	IsAtAll         string   `yaml:"isAtAll"`
}

type statementQiWei struct {
	AlertNotifyUrl  string `yaml:"alertNotifyUrl"`
	MsgType         string `yaml:"msgType"`
	AlertTemplate   string `yaml:"alertTemplate"`
	AlertTitle      string `yaml:"alertTitle"`
	UnknownTemplate string `yaml:"unknownTemplate"`
	UnknownTitle    string `yaml:"unknownTitle"`
}

type statementManifest struct {
	DomainExpireList map[string]string `yaml:"domainExpireList"`
}

type Config struct {
	Global   statementGlobal
	DingDing statementDingDing `yaml:"dingding"`
	QiWei    statementQiWei    `yaml:"qiwei"`
	Domain   []string          `yaml:"domain,flow"`
	Manifest statementManifest `yaml:"manifest"`
}

const FILEPATH string = "config/config.yaml"

var objectConfig *Config

func (c *Config) getConfig() {
	bytes, err := os.ReadFile(FILEPATH)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfigObject() *Config {
	if objectConfig == nil {
		objectConfig = &Config{}
		objectConfig.getConfig()
	}
	return objectConfig
}
