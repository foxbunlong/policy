package driver

import (
	"github.com/oeoen/policy/driver/config"
	"github.com/oeoen/policy/pkg/handler"
	"github.com/oeoen/policy/pkg/police"
	"github.com/oeoen/policy/pkg/tracing"
	"github.com/sirupsen/logrus"
)

type Registry interface {
	Tracer() *tracing.Tracer
	WithConfig(c config.Provider) Registry
	WithLogger(l logrus.FieldLogger) Registry
	Init() error
	police.Registry
	Handler() handler.Provider
}

func CallRegistry(r Registry) {
	r.Tracer()
}
