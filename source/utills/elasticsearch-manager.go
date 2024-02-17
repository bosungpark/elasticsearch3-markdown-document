package utills

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
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
	// create index succeed [200 OK]
	// {
	// 	"acknowledged":true,
	// 	"shards_acknowledged":true,
	// 	"index":"bosung"
	//  }
	_, err := e.client.Indices.Create(name)
	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticManager) IndexDocuments(index string, documents []byte) error {
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
	_, err := e.client.Index(index, bytes.NewReader(documents))
	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticManager) GetDocuments(index string, documentsID string) ([]byte, error) {
	// [200 OK]
	// {
	// 	"_index":"bosung",
	// 	"_id":"frDAto0BFjoKUr_vtenp",
	// 	"_version":1,
	// 	"_seq_no":0,
	// 	"_primary_term":1,
	// 	"found":true,
	// 	"_source":{
	// 	   "message":"테스트가 잘 될까요?"
	// 	}
	//  }
	res, err := e.client.Get(index, documentsID)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
