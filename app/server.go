package app

import (
	"github.com/AnggaPutraa/talk-backend/app/auth"
	"github.com/AnggaPutraa/talk-backend/configs"
	db "github.com/AnggaPutraa/talk-backend/db/sqlc"
	"github.com/AnggaPutraa/talk-backend/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config      configs.Config
	strategy    utils.Strategy
	authHanlder auth.AuthHandler
	router      *gin.Engine
}

func NewServer(config configs.Config, autHander auth.AuthHandler) (*Server, error) {
	server := &Server{
		config: config,
		strategy: utils.NewJWTStrategy(
			config.AccessTokenSecret,
			config.RefreshTokenSecret,
		),
		authHanlder: autHander,
	}
	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	apiGroup := router.Group("/api")
	apiGroup.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	authRoute := apiGroup.Group("/auth")
	authRoute.POST("/register", s.authHanlder.Register)
	authRoute.POST("/login", s.authHanlder.Login)
}

func (s *Server) start(address string) error {
	return s.router.Run(address)
}

func RunServer(config configs.Config, query db.Querier) {
	authService, _ := auth.NewAuthService(config, query)
	authHandler := auth.NewAuthHandler(*authService)
	server, _ := NewServer(config, *authHandler)
	server.start(config.ServerAddress)
}
