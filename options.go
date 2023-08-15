package service

import (
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/options"
)

type clientKey struct{}

// Client to call config service
func Client(c client.Client) options.Option {
	return options.ContextOption(clientKey{}, c)
}

type serviceKey struct{}

// Service to which data load
func Service(s string) options.Option {
	return options.ContextOption(serviceKey{}, s)
}
