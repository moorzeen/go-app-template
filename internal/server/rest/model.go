package rest

import (
	pbAuth "template/internal/service/auth/proto"
	pb "template/internal/service/user/proto"
)

type metaResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type registerResponse struct {
	metaResponse
	User *pb.User `json:"user"`
}

type authResponse struct {
	metaResponse
	Tokens *pbAuth.Tokens `json:"tokens"`
}
