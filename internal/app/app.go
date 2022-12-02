package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go-app-template/internal/broker"
	"go-app-template/internal/server/grpc"
	"go-app-template/internal/server/rest"
	"go-app-template/internal/service/auth"
	"go-app-template/internal/service/user"
	"go-app-template/internal/storage/postgres"
	"go-app-template/internal/storage/redis"
	"go-app-template/pkg"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	config      *config
	userService *user.Service
	authService *auth.Service
	kafkaWriter *kafka.Writer
	grpcServer  *grpc.Server
	restServer  *rest.Server
	shutdownCh  chan os.Signal
}

func Run() error {
	app, err := newApp()
	if err != nil {
		return err
	}

	app.start()
	defer func() {
		err = app.shutdown()
		if err == nil {
			log.Info().Msg("application gracefully stopped")
		} else {
			log.Error().Err(err).Msg("failed to stop application gracefully")
		}
	}()

	select {
	case restErr := <-app.restServer.Error():
		return restErr
	case grpcErr := <-app.grpcServer.Error():
		return grpcErr
	case <-app.shutdownCh:
		return nil
	}
}

func newApp() (*App, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: pkg.TimeLayoutLOG})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Debug().Msg(pkg.AnyPrint(cfg.AppName+" config", cfg))

	st, err := postgres.NewPGStorage(cfg.PG)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	kw := broker.NewKafkaWriter(cfg.Kafka)

	us, err := user.NewUserService(st, kw, cfg.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to init user service: %w", err)
	}

	rc, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to init redis client: %w", err)
	}

	as, err := auth.NewAuthService(st, rc, cfg.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to init auth service: %w", err)
	}

	gs, err := grpc.NewGrpcServer(us, as, cfg.AppGrpcPort)
	if err != nil {
		return nil, fmt.Errorf("failed to init grpc server: %w", err)
	}

	rs := rest.NewRestServer(us, as, cfg.AppRestPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	return &App{
		config:      cfg,
		userService: us,
		authService: as,
		kafkaWriter: kw,
		grpcServer:  gs,
		restServer:  rs,
		shutdownCh:  quit,
	}, nil
}

func (app *App) start() {
	app.grpcServer.Start()
	log.Info().Msgf("grpc server started on port %s", app.config.AppGrpcPort)
	app.restServer.Start()
	log.Info().Msgf("rest server started on port %s", app.config.AppRestPort)
}

func (app *App) shutdown() error {
	app.grpcServer.Stop()
	err := app.restServer.Stop()
	if err != nil {
		return err
	}
	err = app.kafkaWriter.Close()
	if err != nil {
		return err
	}
	return nil
}
