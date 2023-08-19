package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

// 这个文件专门 将日志写入kafka

// 定义一个全局的kafka连接 连接到kafka
var KafkaClient sarama.SyncProducer

type logData struct {
	topic string
	data  string
}

var logDataChan chan *logData

// Kafka连接初始化函数
func KafkaInit(addr []string, size int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	KafkaClient, err = sarama.NewSyncProducer(addr, config)
	logDataChan = make(chan *logData, size)
	if err != nil {
		return
	}
	go SendToKafka()
	return nil
}

func SendToKafka() {
	for {
		select {
		case msgData := <-logDataChan:
			msg := &sarama.ProducerMessage{}
			msg.Topic = msgData.topic
			msg.Value = sarama.StringEncoder(msgData.data)
			KafkaClient.SendMessage(msg)

		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// 用于将日志数据放入到Chan缓冲区
func SendToChan(topic string, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}
