package log

import (
	"context"
	"fmt"
	"sort"

	"github.com/bombsimon/logrusr/v3"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

// LogLevel determines the default log level to be set for any loggers created
// by this package.
var LogLevel = 0

// L either takes a logger from the context (if provided, and if the context
// carries one), or creates a new logger making use of the default logger
// configuration.
//
// example uses:
//
//   - L().Info()
//   - L(ctx).Info()
func L(ctx ...context.Context) logr.Logger {
	switch len(ctx) {
	case 0:
		return New()
	case 1:
		logger, err := logr.FromContext(ctx[0])
		if err != nil {
			return New()
		}

		return logger
	default:
		panic(fmt.Errorf("L: one or zero ctx must be supplied"))
	}
}

// ToContext embeds logger into context `ctx`.
func ToContext(ctx context.Context, logger logr.Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

var fieldsOrder = []string{
	"logger",
	"scenario",
}

// New instantiates a logger making use of default configurations.
func New() logr.Logger {
	logrusLog := logrus.New()
	logrusLog.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		SortingFunc: func(v []string) {
			sort.Strings(v)

			position := 0
			for _, field := range fieldsOrder {
				for idx, item := range v {
					if item == field {
						v[position], v[idx] = v[idx], v[position]
						position++
						break
					}
				}
			}
		},
	})

	logrusLog.SetLevel(logrus.DebugLevel)

	return logrusr.New(logrusLog)
}
