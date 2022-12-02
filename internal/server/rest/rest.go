package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"go-app-template/internal/service/auth"
	pbAuth "go-app-template/internal/service/auth/proto"
	"go-app-template/internal/service/user"
	pb "go-app-template/internal/service/user/proto"
	"net/http"
)

type Server struct {
	server      *http.Server
	userService *user.Service
	authService *auth.Service
	errCh       chan error
}

func NewRestServer(userService *user.Service, authService *auth.Service, port string) *Server {
	g := gin.Default()

	srv := &Server{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: g,
		},
		userService: userService,
		authService: authService,
		errCh:       make(chan error),
	}

	g.POST("/register", srv.register)
	g.GET("/authorize", srv.authorize)

	return srv
}

func (srv *Server) Start() {
	go func() {
		srv.errCh <- srv.server.ListenAndServe()
	}()
}

func (srv *Server) Stop() error {
	return srv.server.Close()
}

func (srv *Server) Error() chan error {
	return srv.errCh
}

func (srv *Server) register(ctx *gin.Context) {
	var req pb.Credentials

	err := jsonpb.Unmarshal(ctx.Request.Body, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(errStatusWithJSON(err))
		return
	}

	registeredUser, err := srv.userService.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.AbortWithStatusJSON(errStatusWithJSON(err))
		return
	}

	response := registerResponse{
		metaResponse: metaResponse{
			StatusCode: http.StatusCreated,
			Message:    "user successfully registered",
		},
		User: registeredUser,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (srv *Server) authorize(ctx *gin.Context) {
	var req pbAuth.Credentials

	err := jsonpb.Unmarshal(ctx.Request.Body, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(errStatusWithJSON(err))
		return
	}

	tokens, err := srv.authService.Authorize(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(errStatusWithJSON(err))
		return
	}

	response := authResponse{
		metaResponse: metaResponse{
			StatusCode: http.StatusOK,
			Message:    "user successfully authorized",
		},
		Tokens: tokens,
	}

	ctx.JSON(http.StatusCreated, response)
}
