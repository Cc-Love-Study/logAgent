package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Cc-Love-Study/LogAgent/conf"
	"github.com/Cc-Love-Study/LogAgent/etcd"
	"github.com/Cc-Love-Study/LogAgent/kafka"
	"github.com/Cc-Love-Study/LogAgent/taillog"
	"gopkg.in/ini.v1"
)

var appConf = conf.NewAppConf()

func init() {
	err := ini.MapTo(appConf, "../conf/config.ini")
	if err != nil {
		fmt.Println("ini load err")
		os.Exit(1)
	}
}

func main() {
	// 0.加载配置文件 在init函数内部
	fmt.Println("ini load success")

	// 1.连接kafka 这里是某kafka的ip端口 开始监听缓冲区
	err := kafka.KafkaInit([]string{appConf.KafkaConf.Address}, appConf.KafkaConf.ChanLen)
	if err != nil {
		fmt.Println("init kafka err")
		return
	}
	fmt.Println("init kafka success")

	// 2. 初始化etcd
	err = etcd.EtcdInit(appConf.EtcdConf.Address, time.Duration(appConf.Timeout)*time.Second)
	if err != nil {
		fmt.Println("init etcd err")
		return
	}
	fmt.Println("init etcd success")
	// 2.1 etcd中拉取配置信息并且监视日志项变化
	confs, err := etcd.GetConf(appConf.EtcdConf.Key)
	if err != nil {
		fmt.Println("conf get failed")
		fmt.Println(err)
		return
	}
	fmt.Println("conf get success")

	// 这里 启用日志追踪管理者 这里已经开始搜集日志 并且发送到缓冲区
	taillog.TailMgrInit(confs)
	fmt.Println("tailMgr init success")

	// 设置 watch 监视配置
	go etcd.WatchConf(appConf.EtcdConf.Key, taillog.SendConfChan())
	fmt.Println("watch conf run success")
	time.Sleep(time.Hour)
}
