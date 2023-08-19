package conf

type LogTransferConf struct {
	KafkaConf `ini:"kafka"`
	EsConf    `ini:"es"`
}

func NewLogTransferConf() *LogTransferConf {
	return new(LogTransferConf)
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type EsConf struct {
	Address string `ini:"address"`
}
