package framework

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	PortNameEnv string
}

type IRouter func(request *Request)

type Server struct {
	engine   *gin.Engine
	server   *http.Server
	database *Db
	config   *Config
	services any
}

func NewServer(config *Config) *Server {
	engine := gin.Default()

	server := &Server{
		engine: engine,
		server: &http.Server{},
		config: config,
	}

	server.loadEnv()

	// engine.Use(yourMiddleware)

	return server
}

func (s *Server) loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (s *Server) Run() {
	s.server.Addr = os.Getenv(s.config.PortNameEnv)
	s.server.Handler = s.engine

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func (s *Server) RegisterRoutes(register []IRouterGroup) {
	for _, group := range register {
		r := s.engine.Group(string(group.GetRouterGroup()))
		for _, route := range group.GetRoutes() {
			route := route // Crear una copia local de la ruta
			r.Handle(route.Method, route.Route, func(c *gin.Context) {
				contexto := c
				route.Controller(
					NewRequest(
						s.GetDatabase(),
						s.GetEngine(),
						s.GetServices(),
						contexto,
					),
				)
			})
		}
	}
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Server) GetEnv(env string) string {
	return os.Getenv(env)
}

func (s *Server) RunDatabase() *Server {
	s.database = NewDatabase(ConfigDatabase{
		User:     s.GetEnv("DB_USER"),
		DbName:   s.GetEnv("DB_NAME"),
		Password: s.GetEnv("DB_PASSWORD"),
		Host:     s.GetEnv("DB_HOST"),
		Port:     s.GetEnv("DB_PORT"),
	})

	return s
}

func (s *Server) RunMigrations() *Server {
	m := NewMigration(s.database.Run())
	err := m.Migrate(s.GetEnv("DIR_MIGRATIONS"))

	if err != nil {
		fmt.Println("error in datgabase migration", err)
	}

	return s
}

func (s *Server) GetDatabase() *Db {
	return s.database
}

func (s *Server) RegisterServices(services any) {
	s.services = services
}

func (s *Server) GetServices() any {
	return s.services
}

func (s *Server) RegisterMiddleware(fn []gin.HandlerFunc) {
	for _, handlerFunc := range fn {
		s.engine.Use(handlerFunc)
	}
}
