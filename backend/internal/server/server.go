package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/furrygem/dia/internal/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Config   *ServerConfig
	listener *net.Listener
	Handler  *httprouter.Router
	Logger   *logrus.Logger
	http.Server
}

func NewServer(config *ServerConfig) (*Server, error) {
	server := &Server{
		Config:  config,
		Handler: httprouter.New(),
		Logger:  logging.GetLogger(),
	}
	return server, nil
}

func (s *Server) getListener() (*net.Listener, error) {
	if s.listener != nil {
		return s.listener, nil
	}

	var err error
	var listener net.Listener
	if s.Config.UseUnixSocketPath {
		listener, err = net.Listen("unix", s.Config.UnixSocketPath)
		if err != nil {
			return nil, err
		}
	} else {
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.BindAddr, s.Config.ListenPort))
		if err != nil {
			return nil, err
		}
	}
	s.listener = &listener

	return s.listener, nil
}

func (s *Server) ListenAndServe() error {
	_, err := s.getListener()
	if err != nil {
		return err
	}
	s.Logger.Infof("Starting listener")
	s.Serve(*s.listener)
	return nil

}
