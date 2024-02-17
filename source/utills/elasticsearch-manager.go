package utills

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticManager struct {
	Client *elasticsearch.Client
}

func (e *ElasticManager) Connect() error {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "mypassword",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true, // 난 로컬에서 띄울거니까 Skip처리
				RootCAs:            nil,  // 필요한 경우 CA 인증서 추가
			},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}
	e.Client = client

	return nil
}

func (e *ElasticManager) CreateIndex(name string) error {
	res, err := e.Client.Indices.Create(name)
	if err != nil {
		return err
	}

	// 예: create index succeed [200 OK] {"acknowledged":true,"shards_acknowledged":true,"index":"bosung"}
	fmt.Println("create index succeed", res)

	return nil
}
