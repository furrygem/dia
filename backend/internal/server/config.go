package server

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	BindAddr                string `yaml:"bind_address"`
	ListenPort              uint16 `yaml:"listen_port"`
	UseUnixSocketPath       bool   `yaml:"use_unix_socket_path"`
	UnixSocketPath          string `yaml:"bind_unix_socket_path"`
	CertificateFilePath     string `yaml:"certificate_file"`
	KeyFilePath             string `yaml:"key_file"`
	PostgresConectionString string `yaml:"postgres_connection_string"`
}

func NewServerConfig() ServerConfig {
	sc := ServerConfig{
		BindAddr:                "0.0.0.0",
		ListenPort:              8000,
		UnixSocketPath:          "",
		UseUnixSocketPath:       false,
		PostgresConectionString: "postgres://postgres:insecure@db:5432/dia",
	}

	return sc
}

func (sc *ServerConfig) FromYAML(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, sc)
	if err != nil {
		return err
	}
	return nil
}

func (sc *ServerConfig) Validate() error {
	if sc.BindAddr != "" {
		if sc.ListenPort == 0 {
			return errors.New("Listen port is not valid")
		}
	}
	return nil
}
