package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

var (
	cfg *Config
)

func InitConfig() {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
		}
	}(file)

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalln(err)
	}
	log.Println("Config is initialized")
	cfg = &config
}

func GetConfig() *Config {
	return cfg
}
