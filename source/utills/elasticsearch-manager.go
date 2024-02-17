package utills

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticManager struct {
	client *elasticsearch.Client
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
	e.client = client

	return nil
}

func (e *ElasticManager) CreateIndex(name string) error {
	res, err := e.client.Indices.Create(name)
	if err != nil {
		return err
	}

	// create index succeed [200 OK]
	// {
	// 	"acknowledged":true,
	// 	"shards_acknowledged":true,
	// 	"index":"bosung"
	//  }
	fmt.Println("create index succeed", res)

	return nil
}

func (e *ElasticManager) IndexDocuments(index string, document []byte) error {
	res, err := e.client.Index(index, bytes.NewReader(document))
	if err != nil {
		return err
	}

	// create index succeed [201 Created]
	// {
	// 	"_index":"bosung",
	// 	"_id":"frDAto0BFjoKUr_vtenp",
	// 	"_version":1,
	// 	"result":"created",
	// 	"_shards":{
	// 	   "total":2,
	// 	   "successful":1,
	// 	   "failed":0
	// 	},
	// 	"_seq_no":0,
	// 	"_primary_term":1
	//  }
	fmt.Println("create index succeed", res)

	return nil
}
