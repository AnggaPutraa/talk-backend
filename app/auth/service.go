package auth

import (
	"context"

	"github.com/AnggaPutraa/talk-backend/configs"
	db "github.com/AnggaPutraa/talk-backend/db/sqlc"
	"github.com/AnggaPutraa/talk-backend/utils"
	"github.com/lib/pq"
)

type AuthService struct {
	config   configs.Config
	query    db.Querier
	strategy utils.Strategy
}

func NewAuthService(config configs.Config, query db.Querier) (*AuthService, error) {
	server := &AuthService{
		config: config,
		query:  query,
		strategy: utils.NewJWTStrategy(
			config.AccessTokenSecret,
			config.RefreshTokenSecret,
		),
	}
	return server, nil
}

func (s *AuthService) CreateUser(c context.Context, request *RegisterRequest) (*TokenResponse, error) {
	hashedPassword := utils.Hash(request.Password)
	createUserParam := db.CreateUserParams{
		Email:          request.Email,
		Username:       request.Username,
		HashedPassword: hashedPassword,
	}
	user, err := s.query.CreateUser(c, createUserParam)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, err
			}
		}
		return nil, err
	}
	accessToken, refreshToken, err := s.strategy.GenerateToken(user.ID, user.Email)
	var response = &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}

func (s *AuthService) LoginUser(c context.Context, request *LoginRequest) (*TokenResponse, error) {
	user, err := s.query.GetUserByEmail(c, request.Email)
	if err != nil {
		return nil, err
	}
	if err := utils.CompareHashed(request.Password, user.HashedPassword); err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := s.strategy.GenerateToken(user.ID, user.Email)
	var response = &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}
