package service

import (
	"context"
	"reflect"

	pb "go.unistack.org/micro-config-service/v5/proto"
	"go.unistack.org/micro/v5/config"
	"go.unistack.org/micro/v5/util/jitter"
	rutil "go.unistack.org/micro/v5/util/reflect"
)

var _ config.Watcher = &serviceWatcher{}

type serviceWatcher struct {
	client  pb.ConfigClient
	service string
	opts    config.Options
	wopts   config.WatchOptions
	done    chan struct{}
	vchan   chan map[string]interface{}
	echan   chan error
}

func (w *serviceWatcher) run() {
	ticker := jitter.NewTicker(w.wopts.MinInterval, w.wopts.MaxInterval)
	defer ticker.Stop()

	src := w.opts.Struct
	if w.wopts.Struct != nil {
		src = w.wopts.Struct
	}

	for {
		select {
		case <-w.done:
			return
		case <-ticker.C:
			dst, err := rutil.Zero(src)
			if err == nil {
				rsp, err := w.client.Load(context.Background(), &pb.LoadRequest{Service: w.service})
				if err == nil {
					err = w.opts.Codec.Unmarshal(rsp.Config, dst)
				}
				if err != nil {
					w.echan <- err
					return
				}
			}
			if err != nil {
				w.echan <- err
				return
			}

			srcmp, err := rutil.StructFieldsMap(src)
			if err != nil {
				w.echan <- err
				return
			}
			dstmp, err := rutil.StructFieldsMap(dst)
			if err != nil {
				w.echan <- err
				return
			}

			for sk, sv := range srcmp {
				if reflect.DeepEqual(dstmp[sk], sv) {
					delete(dstmp, sk)
				}
			}

			if len(dstmp) > 0 {
				w.vchan <- dstmp
				src = dst
			}
		}
	}
}

func (w *serviceWatcher) Next() (map[string]interface{}, error) {
	select {
	case <-w.done:
		break
	case err := <-w.echan:
		return nil, err
	case v, ok := <-w.vchan:
		if !ok {
			break
		}
		return v, nil
	}
	return nil, config.ErrWatcherStopped
}

func (w *serviceWatcher) Stop() error {
	close(w.done)
	return nil
}
