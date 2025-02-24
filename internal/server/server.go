package server

import (
	"log"

	"github.com/MiguelCVA/mhook-backend/internal/api"
	"github.com/gin-gonic/gin"
)

type Server struct {
	port   string
	server *gin.Engine
}

func NewServer() Server {
	return Server{
		port:   "8080",
		server: gin.Default(),
	}
}

func (s *Server) Run() {
	router := api.ConfigRoutes(s.server)

	log.Print(`server is running at port:` + s.port)
	log.Fatal(router.Run(":" + s.port))
}
