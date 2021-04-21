package app

import "go.uber.org/zap"

type cleanupHandler struct {
	logger   *zap.Logger
	cleanups []func()
}

func newCleanupHandler(logger *zap.Logger) *cleanupHandler {
	return &cleanupHandler{
		logger:   logger,
		cleanups: make([]func(), 0),
	}
}

// AddCleanup adds cleanup to cleanup list which is run before application closes
func (ch *cleanupHandler) AddCleanup(cleanup func()) {
	if cleanup == nil {
		ch.logger.Panic("app: cleanupHandler: nil cleanup should not be added")
	}
	ch.cleanups = append(ch.cleanups, cleanup)
}

// AddErrCleanup adds cleanup which returns an error to cleanup list which is run before application closes
func (ch *cleanupHandler) AddErrCleanup(cleanup func() error) {
	if cleanup == nil {
		ch.logger.Panic("app: cleanupHandler: nil cleanup should not be added")
	}
	ch.AddCleanup(func() {
		err := cleanup()
		if err != nil {
			ch.logger.Error("app: cleanupHandler: error when doing cleanup", zap.Error(err))
		}
	})
}

func (ch *cleanupHandler) cleanup() {
	for _, cleanup := range ch.cleanups {
		cleanup()
	}
}
