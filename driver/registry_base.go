package driver

import (
	"github.com/oeoen/policy/driver/config"
	"github.com/oeoen/policy/pkg/handler"
	"github.com/oeoen/policy/pkg/handler/rest"
)

type RegistryBase struct {
	r Registry
	c config.Provider
	h handler.Provider
}

func (b *RegistryBase) with(r Registry) *RegistryBase {
	b.r = r
	return b
}

func (b *RegistryBase) WithConfig(c config.Provider) *RegistryBase {
	b.c = c
	return b
}

func (b *RegistryBase) Handler() handler.Provider {
	if b.c.Service() != "rest" {

	}
	return rest.NewServer(b.r.PoliceManager(), b.c)
}
