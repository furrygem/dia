package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/furrygem/dia/internal/logging"
	"github.com/furrygem/dia/internal/pubkeys"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	http.Server
	Config   *ServerConfig
	listener *net.Listener
	Router   *httprouter.Router
}

func NewServer(config *ServerConfig) (*Server, error) {
	server := &Server{
		Config: config,
		Router: httprouter.New(),
		Server: http.Server{},
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

func (s *Server) addHandlers(pgpool *pgxpool.Pool) {
	logger := logging.GetLogger()
	prefix, pkhs := pubkeys.NewPubKeyHandlers(pgpool).AllRoutes()
	logger.Info(url.Parse)
	for _, handler := range pkhs {
		concatURL, err := url.JoinPath(prefix, handler.Path)
		logger.Infof("Registerig handle %s", concatURL)
		if err != nil {
			logger.Fatal(err.Error())
		}
		s.Router.Handle(handler.Method, concatURL, handler.Handler)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handle http.Handler
	handle = s.Router
	handle = loggingMiddleware(handle)
	handle.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	pgpool, err := pgxpool.New(context.Background(), s.Config.PostgresConectionString)
	if err != nil {
		return err
	}
	logger := logging.GetLogger()
	listener, err := s.getListener()
	s.addHandlers(pgpool)
	s.Server.Handler = s
	if err != nil {
		return err
	}
	logger.Infof("Listening on %s:%d", s.Config.BindAddr, s.Config.ListenPort)
	return s.Serve(*listener)
}
