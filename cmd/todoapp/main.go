package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	core_logger "github.com/swoggoi/todo_list/internal/core/logger"
	core_pgx_pool "github.com/swoggoi/todo_list/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/swoggoi/todo_list/internal/core/transport/http/middleware"
	core_http_server "github.com/swoggoi/todo_list/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/swoggoi/todo_list/internal/features/tasks/repository/postgres"
	tasks_service "github.com/swoggoi/todo_list/internal/features/tasks/service"
	tasks_transport "github.com/swoggoi/todo_list/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/swoggoi/todo_list/internal/features/users/repository/postgres"
	users_service "github.com/swoggoi/todo_list/internal/features/users/service"
	users_transport_http "github.com/swoggoi/todo_list/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg := core_logger.NewConfigMust()
	fmt.Printf("log folder: %s\n", cfg.Folder)

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone ", zap.Any("zone", timeZone))

	logger.Debug("initializing postgres connection pool ")
	pool, err := core_pgx_pool.NewPool(
		core_pgx_pool.NewConfigMust(),
		ctx,
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport.NewTasksService(tasksService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)

	apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion2,
		core_http_middleware.Dummy("api v2 middleware"))
	apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1, apiVersionRouterV2)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
		log.Printf("HTTP server run error: %v", err)
	}
}
