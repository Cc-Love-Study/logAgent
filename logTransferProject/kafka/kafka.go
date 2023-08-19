package kafka

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Cc-Love-Study/logTransfer/es"
	"github.com/Shopify/sarama"
)

func KafKaInit(addr []string, topic string) (err error) {
	// 创建一个连接到消费者的句柄
	consumer, err := sarama.NewConsumer(addr, nil)
	if err != nil {
		fmt.Println("fail to start consumer")
		return
	}
	// 一个话题有好几个分区
	partitionList, err := consumer.Partitions(topic)
	fmt.Println("分区列表:", partitionList)
	// 循环分区
	for patition := range partitionList {
		// 提出分区
		pc, err := consumer.ConsumePartition(topic, int32(patition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("fail start consumer for partition")
			return err
		}

		// 开启线程 服务分区
		fmt.Println("开启分区")
		go func() {
			for msg := range pc.Messages() {
				fmt.Println(string(msg.Value))
				mg := map[string]interface{}{
					"data": string(msg.Value),
				}
				body, err := json.Marshal(mg)
				if err != nil {
					fmt.Println("json err")
					continue
				}
				reader := bytes.NewReader(body)
				es.SendToES(topic, reader)
			}
		}()
	}
	return nil
}
