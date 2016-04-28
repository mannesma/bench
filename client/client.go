package client

type Client interface {
   Get(key string) ([]byte, error)
   Set(key string, value []byte) error
   CreateDir(key string) error
}

func MakeClient(client_type string, server_host string, debug bool) Client {
   switch client_type {
   case "consul":
      return MakeConsulClient(server_host)
      break
   case "etcd":
      return MakeEtcdClient(server_host)
      break
   case "zookeeper":
      client, err := MakeZookeeperClient(server_host, debug)
      if err != nil {
         return nil   
      }
      return client
      break
   }

return nil
}
