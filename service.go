package service

import (
	"context"
	"fmt"

	"github.com/imdario/mergo"
	pb "github.com/unistack-org/micro-config-service/v3/proto"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

var (
	DefaultStructTag = "service"
)

type serviceConfig struct {
	opts    config.Options
	service string
	client  pb.ConfigService
}

func (c *serviceConfig) Options() config.Options {
	return c.opts
}

func (c *serviceConfig) Init(opts ...config.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}

	var cli client.Client
	if c.opts.Context != nil {
		if v, ok := c.opts.Context.Value(serviceKey{}).(string); ok && v != "" {
			c.service = v
		}
		if v, ok := c.opts.Context.Value(clientKey{}).(client.Client); ok {
			cli = v
		}
	}

	if cli == nil {
		return fmt.Errorf("missing Client option")
	}

	if c.service == "" {
		return fmt.Errorf("missing Service option")
	}

	c.client = pb.NewConfigService(c.service, cli)

	return nil
}

func (c *serviceConfig) Load(ctx context.Context) error {
	for _, fn := range c.opts.BeforeLoad {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	rsp, err := c.client.Load(ctx, &pb.LoadRequest{Service: c.service})
	if err != nil && !c.opts.AllowFail {
		return fmt.Errorf("failed to load error config: %w", err)
	}

	src, err := rutil.Zero(c.opts.Struct)
	if err == nil {
		err = c.opts.Codec.Unmarshal(rsp.Config, src)
		if err == nil {
			err = mergo.Merge(c.opts.Struct, src, mergo.WithOverride, mergo.WithTypeCheck, mergo.WithAppendSlice)
		}
	}
	if err != nil && !c.opts.AllowFail {
		return err
	}

	for _, fn := range c.opts.AfterLoad {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	return nil
}

func (c *serviceConfig) Save(ctx context.Context) error {
	for _, fn := range c.opts.BeforeSave {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	for _, fn := range c.opts.AfterSave {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	return nil
}

func (c *serviceConfig) String() string {
	return "service"
}

func NewConfig(opts ...config.Option) config.Config {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}
	return &serviceConfig{opts: options}
}
