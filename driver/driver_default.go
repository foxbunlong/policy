package driver

import (
	"github.com/sirupsen/logrus"

	"github.com/oeoen/policy/driver/config"
)

type DefaultDriver struct {
	c config.Provider
	r Registry
}

func NewDefaultDriver(l logrus.FieldLogger, forcedHTTP bool) Driver {

	c := config.NewViperProvider(l, forcedHTTP)
	var r Registry
	var err error

	if c.DSN() != "" {
		r, err = NewRegistrySqls(c)
	} else {
		l.Fatal("no DSN")
	}

	if err != nil {
		l.WithError(err).Fatal("Unable to instantiate service registry.")
	}

	r.
		WithConfig(c).
		WithLogger(l)

	if err = r.Init(); err != nil {
		l.WithError(err).Fatal("Unable to initialize service registry.")
	}

	return &DefaultDriver{r: r, c: c}
}

func (r *DefaultDriver) Configuration() config.Provider {
	return r.c
}

func (r *DefaultDriver) Registry() Registry {
	return r.r
}

func (r *DefaultDriver) CallRegistry() Driver {
	CallRegistry(r.Registry())
	return r
}
