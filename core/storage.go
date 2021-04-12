package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Storage interface {
	Save(v interface{}) error
	Clear() error
}

type ElasticsearchStorage struct {
	client  *elasticsearch.Client
	context context.Context
	index   string
}

func NewElasticsearchStorage(address []string,
	username string,
	password string,
	index string,
	context context.Context) (*ElasticsearchStorage, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: address,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return &ElasticsearchStorage{es, context, index}, nil
}

func (e ElasticsearchStorage) Save(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index: e.index,
		Body:  bytes.NewReader(data),
	}
	response, err := req.Do(e.context, e.client)
	if err != nil {
		return err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.IsError() {
		return errors.New(response.String())
	}
	return nil
}

func (e ElasticsearchStorage) Clear() error {
	req, err := http.NewRequest("DELETE", "/"+url.PathEscape(e.index), nil)
	if err != nil {
		return err
	}
	req = req.WithContext(e.context)
	response, err := e.client.Transport.Perform(req)
	if err != nil {
		return err
	}
	if (response.StatusCode < 200 || response.StatusCode > 299) && response.StatusCode != 404 {
		return errors.New(response.Status)
	}
	return nil
}
