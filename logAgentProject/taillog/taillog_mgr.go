package taillog

import (
	"fmt"
	"time"

	"github.com/Cc-Love-Study/LogAgent/etcd"
)

var tskMgr *TailTaskMgr

type TailTaskMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func TailMgrInit(confs []*etcd.LogEntry) {
	tskMgr = &TailTaskMgr{
		logEntry:    confs,
		tskMap:      make(map[string]*TailTask, 32),
		newConfChan: make(chan []*etcd.LogEntry), //无缓冲区通道
	}
	for _, conf := range confs {
		// 应该记录tailtask
		task, err := NewTailTask(conf.Path, conf.Topic)
		if err != nil {
			continue
		}
		mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
		tskMgr.tskMap[mk] = task
	}
	go tskMgr.upConf()

}

// 监听配置通道
// 1.配置新增 包括更新
// 2.配置删除
func (t *TailTaskMgr) upConf() {
	for {
		select {
		case res := <-t.newConfChan:
			for _, conf := range res {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mk]
				if ok {
					continue
				} else {
					task, err := NewTailTask(conf.Path, conf.Topic)
					if err != nil {
						continue
					}
					tskMgr.tskMap[mk] = task
				}
			}
			for _, conf := range t.logEntry {
				isDelete := true
				for _, newConf := range res {
					if conf.Path == newConf.Path && conf.Topic == newConf.Topic {
						isDelete = false
					}
				}
				if isDelete {
					mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
					// 停止线程
					t.tskMap[mk].CancelFunc()
					//更新map
					delete(t.tskMap, mk)
					isDelete = true
				}

			}
			// 更新配置表
			t.logEntry = res

			fmt.Println("updata confs")
			fmt.Println(res)
		default:
			time.Sleep(time.Second)
		}
	}
}

func SendConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
