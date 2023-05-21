package services

import (
	"context"
	"net/http"

	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/models"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/pb"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	DB  *mongo.Collection
	Jwt utils.JwtWrapper
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	err := s.DB.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err == nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)

	_, err = s.DB.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	err := s.DB.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.LoginResponse{
				Status: http.StatusNotFound,
				Error:  "User not found",
			}, nil
		}
		return nil, err
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User

	err = s.DB.FindOne(ctx, bson.M{"email": claims.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.ValidateResponse{
				Status: http.StatusNotFound,
				Error:  "User not found",
			}, nil
		}
		return nil, err
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id.Hex(),
	}, nil
}
