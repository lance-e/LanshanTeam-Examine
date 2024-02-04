package serveRegister

import (
	"LanshanTeam-Examine/server/game/utils"
	"context"
	clientv3 "go.etcd.io/etcd/caller/v3"
	"time"
)

type EtcdRegister struct {
	client  *clientv3.Client
	leaseId clientv3.LeaseID
}

func NewEtcdRegister(addr string) (*EtcdRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		utils.GameLogger.Error("newEtcdRegister failed,error:" + err.Error())
		return nil, err
	}
	e := &EtcdRegister{
		client: cli,
	}
	return e, nil
}
func (e *EtcdRegister) ServiceRegister(name string, add string, expire int64) error {
	err := e.CreateLease(expire)
	if err != nil {
		return err
	}
	err = e.BindLease(name, add)
	if err != nil {
		return err
	}
	err = e.KeepAlive()
	return err
}
func (e *EtcdRegister) CreateLease(expire int64) error {
	lease, err := e.client.Grant(context.Background(), expire)
	if err != nil {
		utils.GameLogger.Debug("create lease failed")
		return err
	}
	e.leaseId = lease.ID
	return nil
}
func (e *EtcdRegister) BindLease(key string, value string) error {
	_, err := e.client.Put(context.Background(), key, value, clientv3.WithLease(e.leaseId))
	if err != nil {
		utils.GameLogger.Debug("bind lease failed")
		return err
	}
	return nil
}
func (e *EtcdRegister) KeepAlive() error {
	resp, err := e.client.KeepAlive(context.Background(), e.leaseId)
	if err != nil {
		utils.GameLogger.Debug("keep alive failed")
		return err
	}
	go func(resp <-chan *clientv3.LeaseKeepAliveResponse) {
		for {
			select {
			case <-resp:
				utils.GameLogger.Info("keep alive success")
			}
		}
	}(resp)
	return nil
}
func (e *EtcdRegister) Close() error {
	e.client.Revoke(context.Background(), e.leaseId)
	return e.client.Close()
}
