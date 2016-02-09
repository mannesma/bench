package client

import (
   "errors"
   "fmt"
   "launchpad.net/gozk/zookeeper"
   "time"
)

type ZookeeperClient struct {
   server_host string
   client *zookeeper.Conn
   session <-chan zookeeper.Event
   debug bool
}

func MakeZookeeperClient(server_host string, debug bool) (*ZookeeperClient, error) {
   client := &ZookeeperClient {
      debug: debug,
   }
   if server_host == "" {
      client.server_host = "localhost:2181"
   } else {
      client.server_host = server_host
   }

   conn, session, err := zookeeper.Dial(client.server_host, 5 * time.Second)
   if err != nil {
      fmt.Printf("Couldn't connect: %s\n", err)
      return nil, err
   }

   client.client = conn
   client.session = session

   // Wait for connection.
   event := <-client.session
   fmt.Printf("Got event\n")
   if event.State != zookeeper.STATE_CONNECTED {
      fmt.Printf("Error with connect, %s!\n", event.State)
      return nil, errors.New("Error with connect")
   }

   return client, nil
}

func (z *ZookeeperClient) Get(key string) ([]byte, error) {
   value, _, err := z.client.Get(key)
   if err != nil {
      fmt.Printf("Error with Get: %s\n", err)
      return nil, err
   }
   
   return []byte(value), nil
}

func (z *ZookeeperClient) Set(key string, value []byte) error {
   _, err := z.client.Set(key, string(value), -1)
   if err != nil {
      if e, ok := err.(*zookeeper.Error); ok && e.Code != zookeeper.ZNONODE {
         fmt.Printf("Error with Set: %s\n", err)
         return err
      } else {
         _, err = z.client.Create(key, string(value), 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
         if err != nil {
            fmt.Printf("Error with Create: %s\n", err)
            return err
         }
      }
   }

   return nil
}

func (z *ZookeeperClient) CreateDir(key string) error {
   return z.Set(key, []byte(""))
}
