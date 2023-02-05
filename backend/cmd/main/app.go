package main

import (
	"flag"
	"log"

	"github.com/furrygem/dia/internal/logging"
	"github.com/furrygem/dia/internal/server"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()
}

func main() {
	err := logging.InitLogger("logconfig.yaml")
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
