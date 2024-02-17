package stores

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticStore struct {
	client *elasticsearch.Client
}

func (e *ElasticStore) Connect() error {
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

func (e *ElasticStore) CreateIndex(name string) error {
	// create index succeed [200 OK]
	// {
	// 	"acknowledged":true,
	// 	"shards_acknowledged":true,
	// 	"index":"bosung"
	//  }
	res, err := e.client.Indices.Create(name)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return err
	}

	return nil
}

func (e *ElasticStore) CreateDocuments(index string, documents []byte) error {
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
	res, err := e.client.Index(index, bytes.NewReader(documents))
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
		return err
	}

	return nil
}

func (e *ElasticStore) GetDocument(index string, documentID string) ([]byte, error) {
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
	res, err := e.client.Get(index, documentID)
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

func (e *ElasticStore) SearchDocument(index string, query string) ([]byte, error) {
	// query: `{ "query": { "match_all": {} } }`

	// {
	// 	"took":40,
	// 	"timed_out":false,
	// 	"_shards":{
	// 	   "total":1,
	// 	   "successful":1,
	// 	   "skipped":0,
	// 	   "failed":0
	// 	},
	// 	"hits":{
	// 	   "total":{
	// 		  "value":1,
	// 		  "relation":"eq"
	// 	   },
	// 	   "max_score":1.0,
	// 	   "hits":[
	// 		  {
	// 			 "_index":"bosung",
	// 			 "_id":"frDAto0BFjoKUr_vtenp",
	// 			 "_score":1.0,
	// 			 "_source":{
	// 				"message":"테스트가 잘 될까요?"
	// 			 }
	// 		  }
	// 	   ]
	// 	}
	//  }
	res, err := e.client.Search(
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(strings.NewReader(query)),
	)
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

func (e *ElasticStore) UpdateDocument(index string, documentID string, document string) error {
	// [200 OK]
	// {
	// 	"_index":"bosung",
	// 	"_id":"frDAto0BFjoKUr_vtenp",
	// 	"_version":2,
	// 	"result":"updated",
	// 	"_shards":{
	// 	   "total":2,
	// 	   "successful":1,
	// 	   "failed":0
	// 	},
	// 	"_seq_no":1,
	// 	"_primary_term":1
	//  }
	res, err := e.client.Update(index, documentID, strings.NewReader(document))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return err
	}

	return nil
}

func (e *ElasticStore) DeleteDocument(index string, documentID string) error {
	// [200 OK]
	// {
	// 	"_index":"bosung",
	// 	"_id":"frDAto0BFjoKUr_vtenp",
	// 	"_version":3,
	// 	"result":"deleted",
	// 	"_shards":{
	// 	   "total":2,
	// 	   "successful":1,
	// 	   "failed":0
	// 	},
	// 	"_seq_no":2,
	// 	"_primary_term":1
	// }
	res, err := e.client.Delete(index, documentID)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return err
	}

	return nil
}
