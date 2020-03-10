package driver

import (
	"github.com/oeoen/policy/driver/config"
	"github.com/oeoen/policy/pkg/police"
	"github.com/oeoen/policy/pkg/storage"
	"github.com/oeoen/policy/pkg/storage/sqls"
	"github.com/oeoen/policy/pkg/tracing"
	"github.com/sirupsen/logrus"
)

type RegistrySQLs struct {
	*RegistryBase
	l  logrus.FieldLogger
	c  config.Provider
	t  tracing.Trace
	db storage.Provider
}

func NewRegistrySqls(c config.Provider) (Registry, error) {
	registry := &RegistrySQLs{
		c:            c,
		l:            c.Logger(),
		db:           sqls.NewSQLS(c),
		RegistryBase: new(RegistryBase),
	}
	registry.RegistryBase.with(registry).WithConfig(c)
	if err := registry.Init(); err != nil {
		return nil, err
	}

	return registry, nil
}

func (r *RegistrySQLs) Tracer() *tracing.Tracer {
	return r.Tracer()
}

func (r *RegistrySQLs) WithConfig(c config.Provider) Registry {
	r.c = c
	return r
}

func (r *RegistrySQLs) WithLogger(l logrus.FieldLogger) Registry {
	r.l = l
	return r
}
func (r *RegistrySQLs) Init() error {
	return nil
}
func (r *RegistrySQLs) Provider() storage.Provider {
	return r.db
}
func (r *RegistrySQLs) PoliceManager() police.Manager {
	return r.db
}
func (r *RegistrySQLs) Configuration() police.Configuration {
	return r.c
}
