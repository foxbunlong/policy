package handler

import "github.com/oeoen/policy/pkg/police"

type Provider interface {
	Serve() error
	police.Registry
}
