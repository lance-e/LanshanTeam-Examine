package discovry

import (
	"LanshanTeam-Examine/client/pkg/utils"
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

type EtcdDiscovery struct {
	cli       *clientv3.Client
	serverMap map[string]string
	lock      sync.Mutex
}

func NewServerDiscovery(endpoints []string) (*EtcdDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdDiscovery{
		cli:       client,
		serverMap: make(map[string]string),
	}, nil
}
func (e *EtcdDiscovery) ServerDiscovery(prefix string) error {
	//根据prefix获取注册的服务
	resp, err := e.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, v := range resp.Kvs {
		e.updateServer(string(v.Key), string(v.Value))
	}
	//监听协程，监听prefix变化
	go func() {
		watchchan := e.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
		for watch := range watchchan {
			for _, event := range watch.Events {
				switch event.Type {
				case mvccpb.PUT:
					e.updateServer(string(event.Kv.Key), string(event.Kv.Value))
				case mvccpb.DELETE:
					e.deleteServer(string(event.Kv.Key))
				}
			}
		}

	}()
	return nil
}
func (e *EtcdDiscovery) updateServer(key string, value string) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.serverMap[key] = value
	utils.ClientLogger.Debug("add server :" + key)
}
func (e *EtcdDiscovery) deleteServer(key string) {
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.serverMap, key)
	utils.ClientLogger.Debug("delete server :" + key)
}
