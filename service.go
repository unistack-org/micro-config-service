package service

import (
	"context"
	"fmt"

	"github.com/imdario/mergo"
	pbmicro "github.com/unistack-org/micro-config-service/v3/micro"
	pb "github.com/unistack-org/micro-config-service/v3/proto"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

var DefaultStructTag = "service"

type serviceConfig struct {
	opts    config.Options
	service string
	client  pbmicro.ConfigClient
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
		err := fmt.Errorf("missing client option")
		c.opts.Logger.Error(c.opts.Context, err)
		if !c.opts.AllowFail {
			return err
		}
		return nil
	}

	if c.service == "" {
		err := fmt.Errorf("missing Service option")
		c.opts.Logger.Error(c.opts.Context, err)
		if !c.opts.AllowFail {
			return err
		}
		return nil
	}

	c.client = pbmicro.NewConfigClient(c.service, cli)

	return nil
}

func (c *serviceConfig) Load(ctx context.Context, opts ...config.LoadOption) error {
	if err := config.DefaultBeforeLoad(ctx, c); err != nil {
		return err
	}

	rsp, err := c.client.Load(ctx, &pb.LoadRequest{Service: c.service})
	if err != nil {
		c.opts.Logger.Errorf(ctx, "service load error: %v", err)
		if !c.opts.AllowFail {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return config.DefaultAfterLoad(ctx, c)
	}

	options := config.NewLoadOptions(opts...)
	mopts := []func(*mergo.Config){mergo.WithTypeCheck}
	if options.Override {
		mopts = append(mopts, mergo.WithOverride)
	}
	if options.Append {
		mopts = append(mopts, mergo.WithAppendSlice)
	}

	src, err := rutil.Zero(c.opts.Struct)
	if err == nil {
		err = c.opts.Codec.Unmarshal(rsp.Config, src)
		if err == nil {
			err = mergo.Merge(c.opts.Struct, src, mopts...)
		}
	}

	if err != nil {
		c.opts.Logger.Errorf(ctx, "service load error: %v", err)
		if !c.opts.AllowFail {
			return fmt.Errorf("failed to load config: %w", err)
		}
	}

	if err := config.DefaultAfterLoad(ctx, c); err != nil {
		return err
	}

	return nil
}

func (c *serviceConfig) Save(ctx context.Context, opts ...config.SaveOption) error {
	if err := config.DefaultBeforeSave(ctx, c); err != nil {
		return err
	}

	buf, err := c.opts.Codec.Marshal(c.opts.Struct)
	if err == nil {
		_, err = c.client.Save(ctx, &pb.SaveRequest{Service: c.service, Config: buf})
	}
	if err != nil {
		c.opts.Logger.Errorf(ctx, "service save error: %v", err)
		if !c.opts.AllowFail {
			return fmt.Errorf("failed to save config: %w", err)
		}
	}

	if err := config.DefaultAfterSave(ctx, c); err != nil {
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

func (c *serviceConfig) Watch(ctx context.Context, opts ...config.WatchOption) (config.Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}

func NewConfig(opts ...config.Option) config.Config {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}
	return &serviceConfig{opts: options}
}
