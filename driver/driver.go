package driver

import (
	"github.com/oeoen/policy/driver/config"
)

type Driver interface {
	Configuration() config.Provider
	Registry() Registry
	CallRegistry() Driver
}
