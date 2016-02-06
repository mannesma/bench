package client

import (
   "fmt"
)

type Client interface {
   Get(key string) ([]byte, error)
   Set(key string, value []byte) error
}

func MakeClient(client_type string, server_host string) Client {
   switch client_type {
   case "consul":
      return MakeConsulClient(server_host)
      break
   case "etcd":
      // if server_host == "" {
      //   server_host = "localhost:2379"
      // }
      // client.base_url = fmt.Sprintf("http://%s/v2/keys", server_host)
      fmt.Printf("etcd not implemented")
      return nil
      break
   case "zookeeper":
      fmt.Printf("zookeeper not implemented")
      return nil
      break
   }

return nil
}
