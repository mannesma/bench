package main

import (
	"errors"
	"fmt"
   "github.com/mannesma/bench"
	"os"
	"time"
)

import "launchpad.net/gozk/zookeeper"

func MakeTestHttpClient(config *bench.Config) *TestHttpClient {
   var server_host string = ""

   client := &TestHttpClient {}
   if config.ServerHost != "" {
      server_host = config.ServerHost
   }
   switch config.ServerType {
   case "consul":
      if server_host == "" {
         server_host = "localhost:8500"
      }
      client.base_url = fmt.Sprintf("http://%s/v1/kv", server_host)
      break
   case "etcd":
      if server_host == "" {
         server_host = "localhost:2379"
      }
      client.base_url = fmt.Sprintf("http://%s/v2/keys", server_host)
      break
   }
   client.client = &http.Client{}

   return client
}

type TestZookeeperClient struct {
   server_host string
   client *zookeeper.Conn
   session <-chan zookeeper.Event
   debug bool
}

func MakeTestZookeeperClient(config *bench.Config) (*TestZookeeperClient, error) {
   client := &TestZookeeperClient {
      debug: config.Debug,
   }
   if config.ServerHost == "" {
      client.server_host = "localhost:2181"
   } else {
      client.server_host = config.ServerHost
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

func (tzc *TestZookeeperClient) Get(key string) ([]byte, error) {
   value, _, err := tzc.client.Get(key)
   if err != nil {
      fmt.Printf("Error with Get: %s\n", err)
      return nil, err
   }
   
   return []byte(value), nil
}

func (tzc *TestZookeeperClient) Set(key string, value []byte) error {
   _, err := tzc.client.Set(key, string(value), -1)
   if err != nil {
      if e, ok := err.(*zookeeper.Error); ok && e.Code != zookeeper.ZNONODE {
         fmt.Printf("Error with Set: %s\n", err)
         return err
      } else {
         _, err = tzc.client.Create(key, string(value), 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
         if err != nil {
            fmt.Printf("Error wit Create: %s\n", err)
            return err
         }
      }
   }

   return nil
}

func main() {
   config := bench.MakeConfigFromCmdline()
   t := MakeSuite(config)   

   if t == nil {
      fmt.Printf("Error creating Suite!\n")
      os.Exit(1)
   }

   if t.Config.Setup {
      err := t.SetupTest()
      if err != nil {
         fmt.Printf("Error with setup: %s\n", err)
      }
   } else {
      t.run_test()
      t.ReadPerf.Print(true)
   }
}
