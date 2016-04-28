package client

import (
   "fmt"
)

type EtcdClient struct {
   client *HttpClient
}

func MakeEtcdClient(server_host string) *EtcdClient {
   if server_host == "" {
         server_host = "localhost:2379"
   }
   base_url := fmt.Sprintf("http://%s/v2/keys", server_host)
   client := &EtcdClient {
      client: MakeHttpClient(base_url),
   } 

   return client
}

func (e *EtcdClient) Get(key string) ([]byte, error) {
   return e.client.Get(key)
}

func (e *EtcdClient) Set(key string, value []byte) error {
   body := fmt.Sprintf("value=%s", value)
   return e.client.Put(key, []byte(body))
}

func (e *EtcdClient) CreateDir(key string) error {
   return e.client.Put(key, []byte("dir=true"))
}
