package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Cc-Love-Study/logTransfer/conf"
	"github.com/Cc-Love-Study/logTransfer/es"
	"github.com/Cc-Love-Study/logTransfer/kafka"
	"gopkg.in/ini.v1"
)

var LogConf = conf.NewLogTransferConf()

func init() {
	err := ini.MapTo(LogConf, "./conf/cfg.ini")
	if err != nil {
		fmt.Println("ini load err")
		os.Exit(1)
	}
	fmt.Println("ini load success")
}

// 将日志从kafka取出 发往ES
func main() {
	// 0.加载配置
	// 1.初始化kafka
	err := kafka.KafKaInit([]string{LogConf.KafkaConf.Address}, LogConf.Topic)
	if err != nil {
		fmt.Println("kafka init err")
		return
	}
	fmt.Println("kafka init success")
	// 2.初始化es
	err = es.ESInit(LogConf.EsConf.Address)
	if err != nil {
		fmt.Println("es init err")
		return
	}
	fmt.Println("Es init success")
	time.Sleep(time.Hour)
}
