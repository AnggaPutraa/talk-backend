package app

import (
	"github.com/AnggaPutraa/talk-backend/app/auth"
	"github.com/AnggaPutraa/talk-backend/app/ws"
	"github.com/AnggaPutraa/talk-backend/configs"
	db "github.com/AnggaPutraa/talk-backend/db/sqlc"
	"github.com/AnggaPutraa/talk-backend/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config           configs.Config
	strategy         utils.Strategy
	authHanlder      auth.AuthHandler
	webSocketHandler ws.WebSocketHandler
	router           *gin.Engine
}

func NewServer(
	config configs.Config,
	autHander auth.AuthHandler,
	webSocketHandler ws.WebSocketHandler,
) (*Server, error) {
	server := &Server{
		config: config,
		strategy: utils.NewJWTStrategy(
			config.AccessTokenSecret,
			config.RefreshTokenSecret,
		),
		authHanlder:      autHander,
		webSocketHandler: webSocketHandler,
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
	webSocketRoute := apiGroup.Group("/ws/room")
	webSocketRoute.GET("/", s.webSocketHandler.GetRooms)
	webSocketRoute.POST("/", s.webSocketHandler.CreateRoom)
	webSocketRoute.GET("/join/:id", s.webSocketHandler.JoinRoom)
	webSocketRoute.GET("/:id/client", s.webSocketHandler.GetClientsByRoomId)
	s.router = router
}

func (s *Server) start(address string) error {
	return s.router.Run(address)
}

func RunServer(config configs.Config, query db.Querier) {
	authService, _ := auth.NewAuthService(config, query)
	authHandler := auth.NewAuthHandler(*authService)
	hub, _ := ws.NewHub()
	webSocketHandler := ws.NewWebSocketHandler(*hub)
	server, _ := NewServer(config, *authHandler, *webSocketHandler)
	go hub.Run()
	server.start(config.ServerAddress)
}
