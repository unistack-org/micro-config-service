package service // import "go.unistack.org/micro-config-service/v4"

import (
	"context"
	"fmt"

	"dario.cat/mergo"
	pbmicro "go.unistack.org/micro-config-service/v4/micro"
	pb "go.unistack.org/micro-config-service/v4/proto"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/config"
	"go.unistack.org/micro/v4/options"
	rutil "go.unistack.org/micro/v4/util/reflect"
)

var _ config.Config = &serviceConfig{}

var DefaultStructTag = "service"

type serviceConfig struct {
	opts    config.Options
	service string
	client  pb.ConfigServiceClient
}

func (c *serviceConfig) Options() config.Options {
	return c.opts
}

func (c *serviceConfig) Init(opts ...options.Option) error {
	if err := config.DefaultBeforeInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
		return err
	}

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
		err := fmt.Errorf("missing client option")
		if !c.opts.AllowFail {
			return err
		}

		if err := config.DefaultAfterInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
			return err
		}

		return nil
	}

	if c.service == "" {
		err := fmt.Errorf("missing Service option")
		if !c.opts.AllowFail {
			return err
		}

		if err := config.DefaultAfterInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
			return err
		}

		return nil
	}

	c.client = pbmicro.NewConfigServiceClient(c.service, cli)

	if err := config.DefaultAfterInit(c.opts.Context, c); err != nil && !c.opts.AllowFail {
		return err
	}

	return nil
}

func (c *serviceConfig) Load(ctx context.Context, opts ...options.Option) error {
	if err := config.DefaultBeforeLoad(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	rsp, err := c.client.Load(ctx, &pb.LoadRequest{Service: c.service})
	if err != nil {
		if !c.opts.AllowFail {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err = config.DefaultAfterLoad(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}

		return nil
	}

	options := config.NewLoadOptions(opts...)
	mopts := []func(*mergo.Config){mergo.WithTypeCheck}
	if options.Override {
		mopts = append(mopts, mergo.WithOverride)
	}
	if options.Append {
		mopts = append(mopts, mergo.WithAppendSlice)
	}

	dst := c.opts.Struct
	if options.Struct != nil {
		dst = options.Struct
	}

	src, err := rutil.Zero(dst)
	if err == nil {
		err = c.opts.Codec.Unmarshal(rsp.Config, src)
		if err == nil {
			err = mergo.Merge(dst, src, mopts...)
		}
	}

	if err != nil && !c.opts.AllowFail {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := config.DefaultAfterLoad(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	return nil
}

func (c *serviceConfig) Save(ctx context.Context, opts ...options.Option) error {
	if err := config.DefaultBeforeSave(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	options := config.NewSaveOptions(opts...)

	dst := c.opts.Struct
	if options.Struct != nil {
		dst = options.Struct
	}

	buf, err := c.opts.Codec.Marshal(dst)
	if err == nil {
		_, err = c.client.Save(ctx, &pb.SaveRequest{Service: c.service, Config: buf})
	}
	if err != nil && !c.opts.AllowFail {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := config.DefaultAfterSave(ctx, c); err != nil && !c.opts.AllowFail {
		return err
	}

	return nil
}

func (c *serviceConfig) String() string {
	return "service"
}

func (c *serviceConfig) Name() string {
	return c.opts.Name
}

func (c *serviceConfig) Watch(ctx context.Context, opts ...options.Option) (config.Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}

func NewConfig(opts ...options.Option) *serviceConfig {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}
	return &serviceConfig{opts: options}
}
