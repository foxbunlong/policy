package police

import "net/http"

type Validator struct {
	cl *http.Client
	c  Configuration
}

func NewValidator(conf Configuration) *Validator {
	return &Validator{
		cl: http.DefaultClient,
		c:  conf,
	}
}

func NewValidatorWithClient(conf Configuration, client *http.Client) *Validator {
	return &Validator{
		cl: client,
		c:  conf,
	}
}
