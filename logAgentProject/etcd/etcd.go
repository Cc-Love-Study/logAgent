package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 一个路径 一个发送到的话题
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

var etcdClient *clientv3.Client

func EtcdInit(addr string, t time.Duration) (err error) {
	etcdClient, err = clientv3.New(
		clientv3.Config{
			Endpoints:   []string{addr},
			DialTimeout: t,
		})
	if err != nil {
		return
	}
	return nil
}

// key：value 这里的v是一个json格式的字符串 里面有path  和 topic
func GetConf(key string) (confs []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	res, err := etcdClient.Get(ctx, key)
	cancel()
	if err != nil {
		return
	}
	for _, ev := range res.Kvs {
		fmt.Println(string(ev.Key), string(ev.Value))
		err = json.Unmarshal(ev.Value, &confs)
		if err != nil {
			return
		}
	}
	return
}

func WatchConf(key string, newChan chan<- []*LogEntry) {
	// 设置监视
	ch := etcdClient.Watch(context.Background(), key)

	// 每次有数据更新 就返回一个值到ch
	for wresp := range ch {
		for _, evt := range wresp.Events {
			// 拿到最新的配置 传输到taillog.Mgr
			var confs = make([]*LogEntry, 0)
			// 删除操作
			if evt.Type == clientv3.EventTypeDelete {
				newChan <- confs
				continue
			}
			if evt.Type == clientv3.EventTypePut {
				err := json.Unmarshal(evt.Kv.Value, &confs)
				if err != nil {
					fmt.Println("err json")
					continue
				}
				newChan <- confs
			}

		}
	}
}
