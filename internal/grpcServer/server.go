package grpcServer

import (
	"context"
	"time"

	"github.com/google/uuid"

	auth "github.com/fpmi-hci-2023/project13b-auth/api/auth/v1/gen"
	"github.com/fpmi-hci-2023/project13b-auth/config"
	"github.com/fpmi-hci-2023/project13b-auth/internal/db"
	"github.com/fpmi-hci-2023/project13b-auth/internal/model"
	"github.com/fpmi-hci-2023/project13b-auth/internal/token"
	"github.com/fpmi-hci-2023/project13b-auth/pkg/logger"
)

type GRPCServer struct {
	auth.UnimplementedAuthServiceServer
	log logger.Logger
	db  *db.DB
}

func (s *GRPCServer) Register(_ context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	user := model.User{
		ID:       uuid.New().String(),
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	err := s.db.AddUser(&user)
	if err != nil {
		s.log.Debug(err)

		return &auth.RegisterResponse{}, err
	}

	accessToken, err := token.NewEncryptedToken(user.ID, user.Email, user.Username, config.GlobalConfig.TTL.Access, config.GlobalConfig.AES)
	if err != nil {
		s.log.Error(err, "error while creating access token")

		return &auth.RegisterResponse{}, err
	}

	refreshToken, err := token.NewEncryptedToken(user.ID, user.Email, user.Username, config.GlobalConfig.TTL.Refresh, config.GlobalConfig.AES)
	if err != nil {
		s.log.Error(err, "error while creating refresh token")

		return &auth.RegisterResponse{}, err
	}

	return &auth.RegisterResponse{
		Id:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken.EncryptedToken,
		RefreshToken: refreshToken.EncryptedToken,
	}, nil
}

func (s *GRPCServer) Login(_ context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	var user model.User

	user.Username = req.Username
	user.Password = req.Password

	err := s.db.GetUser(&user)
	if err != nil {
		s.log.Debug(err)

		return &auth.LoginResponse{}, err
	}

	accessToken, err := token.NewEncryptedToken(user.ID, user.Email, user.Username, config.GlobalConfig.TTL.Access, config.GlobalConfig.AES)
	if err != nil {
		s.log.Error(err, "error while creating access token")

		return &auth.LoginResponse{}, err
	}

	refreshToken, err := token.NewEncryptedToken(user.ID, user.Email, user.Username, config.GlobalConfig.TTL.Refresh, config.GlobalConfig.AES)
	if err != nil {
		s.log.Error(err, "error while creating refresh token")

		return &auth.LoginResponse{}, err
	}

	return &auth.LoginResponse{
		Id:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		LastLogin:    user.LastLogin.String(),
		AccessToken:  accessToken.EncryptedToken,
		RefreshToken: refreshToken.EncryptedToken,
	}, nil
}

// Validate checks access and refresh tokens.
func (s *GRPCServer) Validate(_ context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	accessToken := token.EncryptedJWT{
		EncryptedToken: req.AccessToken,
		Key:            config.GlobalConfig.AES,
	}
	refreshToken := token.EncryptedJWT{
		EncryptedToken: req.RefreshToken,
		Key:            config.GlobalConfig.AES,
	}

	if err := accessToken.Check(); err != nil {
		s.log.Debug("access token is invalid")

		if err = refreshToken.Check(); err != nil {
			s.log.Debug("refresh token is invalid")

			return &auth.ValidateResponse{}, err
		}

		refreshTokenClaims, err := refreshToken.Parse()
		if err != nil {
			s.log.Error(err, "error while parsing refresh token")

			return &auth.ValidateResponse{}, err
		}

		updatedAccessToken, err := token.NewEncryptedToken(
			(*refreshTokenClaims)["id"].(string),
			(*refreshTokenClaims)["email"].(string),
			(*refreshTokenClaims)["username"].(string),
			config.GlobalConfig.TTL.Access,
			config.GlobalConfig.AES,
		)
		if err != nil {
			s.log.Error(err, "error while creating new access token")

			return &auth.ValidateResponse{}, err
		}

		updatedRefreshToken, err := token.NewEncryptedToken(
			(*refreshTokenClaims)["id"].(string),
			(*refreshTokenClaims)["email"].(string),
			(*refreshTokenClaims)["username"].(string),
			config.GlobalConfig.TTL.Refresh,
			config.GlobalConfig.AES,
		)
		if err != nil {
			s.log.Error(err, "error while creating new refresh token")

			return &auth.ValidateResponse{}, err
		}

		return &auth.ValidateResponse{
			TokenStatus:  auth.ValidateResponse_UPDATE,
			Id:           (*refreshTokenClaims)["id"].(string),
			Username:     (*refreshTokenClaims)["username"].(string),
			Email:        (*refreshTokenClaims)["email"].(string),
			AccessToken:  updatedAccessToken.EncryptedToken,
			RefreshToken: updatedRefreshToken.EncryptedToken,
		}, nil
	}

	accessTokenClaims, err := accessToken.Parse()
	if err != nil {
		s.log.Error(err, "error while parsing access token")

		return &auth.ValidateResponse{}, err
	}

	return &auth.ValidateResponse{
		TokenStatus:  auth.ValidateResponse_OK,
		Id:           (*accessTokenClaims)["id"].(string),
		Username:     (*accessTokenClaims)["username"].(string),
		Email:        (*accessTokenClaims)["email"].(string),
		AccessToken:  accessToken.EncryptedToken,
		RefreshToken: refreshToken.EncryptedToken,
	}, nil
}

func (s *GRPCServer) Info(_ context.Context, req *auth.ValidateRequest) (*auth.InfoResponse, error) {
	accessToken := token.EncryptedJWT{
		EncryptedToken: req.AccessToken,
		Key:            config.GlobalConfig.AES,
	}
	refreshToken := token.EncryptedJWT{
		EncryptedToken: req.RefreshToken,
		Key:            config.GlobalConfig.AES,
	}

	if err := accessToken.Check(); err != nil {
		s.log.Debug("access token is invalid")

		if err = refreshToken.Check(); err != nil {
			s.log.Debug("refresh token is invalid")

			return &auth.InfoResponse{}, err
		}

		refreshTokenClaims, err := refreshToken.Parse()
		if err != nil {
			s.log.Error(err, "error while parsing refresh token")

			return &auth.InfoResponse{}, err
		}

		return &auth.InfoResponse{
			Id:        (*refreshTokenClaims)["id"].(string),
			Username:  (*refreshTokenClaims)["email"].(string),
			Email:     (*refreshTokenClaims)["username"].(string),
			LastLogin: time.Unix(int64((*refreshTokenClaims)["exp"].(float64))-config.GlobalConfig.TTL.Refresh, 0).String(),
		}, nil
	}

	accessTokenClaims, err := accessToken.Parse()
	if err != nil {
		s.log.Error(err, "error while parsing access token")

		return &auth.InfoResponse{}, err
	}

	return &auth.InfoResponse{
		Id:        (*accessTokenClaims)["id"].(string),
		Username:  (*accessTokenClaims)["username"].(string),
		Email:     (*accessTokenClaims)["email"].(string),
		LastLogin: time.Unix(int64((*accessTokenClaims)["exp"].(float64))-config.GlobalConfig.TTL.Access, 0).String(),
	}, nil
}
