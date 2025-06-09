package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type CommonConfiguration struct {
	Port string `default:":8090"`
}

var CommonConfig CommonConfiguration

func LoadConfig() {
	if err := envconfig.Process("", &CommonConfig); err != nil {
		log.Panicf("Failed to load CommonConfig. Error: %v", err)
	}
	log.Printf("CommonConfig : %+v", CommonConfig)
}
