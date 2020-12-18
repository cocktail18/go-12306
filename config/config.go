package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	Cookie     string   `json:"cookie"`
	Train      []string `json:"train"`
	Seats      []string `json:"seats"`
	Passengers []string `json:"passengers"`
	Form       string   `json:"form"`
	To         string   `json:"to"`
	Time       string   `json:"time"`
	Mail       []string `json:"mail"`
	Fmail      string   `json:"fmaul"`
	Fpassword  string   `json:"fpassword"`
	Host       string   `json:"host"`
	Port       string   `json:"port"`
}

func (c *Config) GetConf() *Config {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
	}
	return c
}
