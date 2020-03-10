package police

import "github.com/oeoen/policy/driver/config"

type Registry interface {
	PoliceManager() Manager
	Configuration() Configuration
	// Validator() *Validator
}

type Configuration interface {
	config.Provider
}
