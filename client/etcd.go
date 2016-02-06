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

func (c *EtcdClient) Get(key string) ([]byte, error) {
   return c.client.Get(key)
}

func (c *EtcdClient) Set(key string, value []byte) error {
   return c.client.Put(key, value)
}
