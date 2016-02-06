package client

import (
   "fmt"
)

type ConsulClient struct {
   client *HttpClient
}

func MakeConsulClient(server_host string) *ConsulClient {
   if server_host == "" {
         server_host = "localhost:8500"
   }
   base_url := fmt.Sprintf("http://%s/v1/kv", server_host)
   client := &ConsulClient {
      client: MakeHttpClient(base_url),
   } 

   return client
}

func (c *ConsulClient) Get(key string) ([]byte, error) {
   return c.client.Get(key)
}

func (c *ConsulClient) Set(key string, value []byte) error {
   return c.client.Put(key, value)
}
