package main

import (
	"flag"
	"log"

	"github.com/furrygem/dia/internal/logging"
	"github.com/furrygem/dia/internal/server"
)

var configPath string
var loggingConfigPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.StringVar(&loggingConfigPath, "logging-config", "", "Path to logging configuration file")
	flag.Parse()
}

func main() {
	err := logging.InitLogger(loggingConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	serverConfig := server.NewServerConfig()
	if configPath != "" {
		err := serverConfig.FromYAML(configPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	s, err := server.NewServer(&serverConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	s.ListenAndServe()

}
