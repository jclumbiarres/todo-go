package server

import (
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Server struct {
	Config      Config
	router      *http.ServeMux
	middlewares []func(http.Handler) http.Handler
	DB          *gorm.DB
}

type Config struct {
	listenaddress string
	id            string
	name          string
}

func NewConfig() Config {
	return Config{
		listenaddress: ":3000",
		id:            "s01",
		name:          "server01",
	}
}

func (c Config) SetListenAddress(listen string) Config {
	c.listenaddress = listen
	return c
}

func (c Config) SetID(id string) Config {
	c.id = id
	return c
}

func (c Config) SetName(name string) Config {
	c.name = name
	return c
}

func NewServer(config Config, db *gorm.DB) *Server {
	return &Server{
		Config:      config,
		router:      http.NewServeMux(),
		middlewares: []func(http.Handler) http.Handler{},
		DB:          db,
	}
}

func (s *Server) Use(middleware func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *Server) Route(path string, handler http.HandlerFunc) {
	finalHandler := handler
	for _, middleware := range s.middlewares {
		finalHandler = middleware(finalHandler).(http.HandlerFunc)
	}
	s.router.HandleFunc(path, finalHandler)
}

func (s *Server) Handle(path string, handler http.Handler) {
	finalHandler := handler
	for _, middleware := range s.middlewares {
		finalHandler = middleware(finalHandler).(http.HandlerFunc)
	}
	s.router.Handle(path, handler)
}

func (s *Server) Start() {
	log.Printf("\tID: %s\t\tName: %s\tListen: %s", s.Config.id, s.Config.name, s.Config.listenaddress)
	err := http.ListenAndServe(s.Config.listenaddress, s.router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
