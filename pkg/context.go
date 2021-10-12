package pkg

import (
	"context"
	"go.uber.org/zap"
)

type Context struct {
	Context context.Context
	Logger *zap.SugaredLogger
}
