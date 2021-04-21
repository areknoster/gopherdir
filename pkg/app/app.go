package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"

	"go.uber.org/zap/zapcore"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Config is used to configure application
type Config struct {
	HTTPHostAddress         string `default:"0.0.0.0:80" split_words:"true"`
	LogErrorsWithStacktrace bool   `default:"false" split_words:"true"`
}

// App is responsible for initializing and running application
type App struct {
	Logger     *zap.Logger
	Router     chi.Router
	httpServer *http.Server
	config     Config
	*cleanupHandler
}

// NewApp initializes application
func NewApp(logger *zap.Logger, cfg Config) *App {
	if !cfg.LogErrorsWithStacktrace {
		logger.WithOptions(zap.AddStacktrace(zapcore.PanicLevel))
	}

	router := chi.NewRouter().
		With(
			middleware.DefaultLogger,
			middleware.Recoverer,
		)

	return &App{
		Logger:         logger,
		Router:         router,
		config:         cfg,
		cleanupHandler: newCleanupHandler(logger),
	}
}

// Run runs application server and other services
func (a *App) Run() {
	defer a.cleanup()
	a.Logger.Info("Application started running",
		zap.String("address", a.config.HTTPHostAddress))
	a.httpServer = &http.Server{
		Addr:    a.config.HTTPHostAddress,
		Handler: a.Router,
	}

	if err := a.httpServer.ListenAndServe(); err != nil {
		a.Logger.Fatal("app: error while http listening", zap.Error(err))
	}
}

// Build is used to initialize application's dependencies
func (a *App) Build(initializer func() error) {
	err := initializer()
	if err != nil {
		a.Logger.Error("app: could not initialize application", zap.Error(err))
		os.Exit(1)
	}
}
