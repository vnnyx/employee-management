package server

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	configapi "github.com/vnnyx/employee-management/config/api"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/pkg/apperror"
	"go.uber.org/zap"
)

type Server struct {
	Config configapi.Config
	Logger *zap.SugaredLogger
	Fiber  *fiber.App
	DB     *pgxpool.Pool
}

func NewServer(config configapi.Config, logger *zap.SugaredLogger, db *pgxpool.Pool) *Server {
	var fiberConfig fiber.Config
	fiberConfig.ErrorHandler = apperror.HTTPHandleError
	fiberConfig.AppName = config.App.Name
	fiberConfig.DisableStartupMessage = true
	fiberConfig.JSONEncoder = json.Marshal
	fiberConfig.JSONDecoder = json.Unmarshal
	fiberConfig.ReadBufferSize = 16 * 1024
	fiberConfig.BodyLimit = 100 * 1024 * 1024
	fiberConfig.EnableSplittingOnParsers = true

	return &Server{
		Config: config,
		Logger: logger,
		Fiber:  fiber.New(fiberConfig),
		DB:     db,
	}
}

func (s *Server) Run() error {
	// Request Logger Middleware

	if s.Config.App.Env == "local" {
		config := logger.ConfigDefault
		config.Format = "[${time}] ${status} ${method} ${originalURL}\n"
		config.CustomTags = map[string]logger.LogFunc{
			"originalURL": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				decodedValue, err := url.QueryUnescape(c.OriginalURL())
				if err != nil {
					return 0, errors.Wrap(err, constants.ErrWrapUrlQueryUnescape)
				}

				return output.WriteString(decodedValue)
			},
		}
		s.Fiber.Use(logger.New(config))
	}

	// Trace Middleware
	if *s.Config.Observability.Enable {
		s.Fiber.Use(
			skip.New(
				otelfiber.Middleware(otelfiber.WithoutMetrics(true)),
				func(c *fiber.Ctx) bool {
					return c.Method() == http.MethodOptions ||
						strings.Contains(c.Path(), "/health/check") ||
						strings.Contains(c.Path(), "/public")
				},
			),
		)
	}

	// Recover Middleware
	s.Fiber.Use(recover.New(recover.Config{EnableStackTrace: true, StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
		c.Locals("recoveredStack", debug.Stack())
	}}))

	s.Fiber.Use(compress.New())

	// Swagger Handler
	s.Fiber.Get("/swagger/*", swagger.HandlerDefault)

	// Map App Handlers
	err := s.MapHandlers()
	if err != nil {
		return err
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quit
		s.Fiber.Shutdown()
	}()

	// Run Fiber
	s.Logger.Infof("App started")
	s.Logger.Infof("Listening at :%d", s.Config.App.Port)
	return s.Fiber.Listen(fmt.Sprintf(":%s", strconv.FormatInt(s.Config.App.Port, 10)))
}
