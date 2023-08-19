package es

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
)

var client *elastic.Client

func ESInit(address string) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	cfg := elastic.Config{
		Addresses: []string{
			address,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}
	client, err = elastic.NewClient(cfg)
	if err != nil {
		return
	}
	return nil
}

func SendToES(indexStr string, data io.Reader) (err error) {
	// 编辑索引 类型 数据    index      data
	//                     数据表名   数据内容
	_, err = client.Index(indexStr, data)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("发送成功")
	// fmt.Println(put.Id, put.Index, put.Type)
	return
}
