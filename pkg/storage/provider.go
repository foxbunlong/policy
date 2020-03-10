package storage

import (
	"github.com/oeoen/policy/driver/config"
	"github.com/oeoen/policy/pkg/police"
)

type Provider interface {
	DBInit(c config.Provider) error
	DBDefer() error
	police.Manager
}
