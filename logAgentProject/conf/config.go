package conf

type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

func NewAppConf() *AppConf {
	return new(AppConf)
}

type KafkaConf struct {
	Address string `ini:"address"`
	ChanLen int    `ini:"chanmaxlen"`
}

type EtcdConf struct {
	Key     string `ini:"key"`
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
}

//--------------------------------------------------------
type TaillogConf struct {
	FileName string `ini:"filename"`
}
