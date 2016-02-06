package main

import (
	"bytes"
	"errors"
	"fmt"
   "github.com/mannesma/bench"
	"github.com/mannesma/bench/perf"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

import "launchpad.net/gozk/zookeeper"

type TestHttpClient struct {
   base_url string
   client *http.Client
}

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

func (thc *TestHttpClient) Get(key string) ([]byte, error) {
   resp, err := thc.client.Get(thc.base_url + key)
   if err != nil {
      fmt.Printf("Error with get: %s\n", err)
      return nil, err
   }
   defer resp.Body.Close()
   value, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      fmt.Printf("Error with Get: %s\n", err)
      return nil, err
   }

   if resp.StatusCode != 200 {
      fmt.Printf("Error: key = %s (%d) %s\n", key, resp.StatusCode, value)
      valstr := string(value)
      return nil, errors.New(valstr)
   }
   
   return value, nil
}

func (thc *TestHttpClient) Set(key string, value []byte) error {
   br := bytes.NewReader(value)
   req, err := http.NewRequest("PUT", thc.base_url + key, br)
   if err != nil {
      fmt.Printf("Error with Set request: %s\n", err)
      return err
   }
   resp, err := thc.client.Do(req)
   if err != nil {
      fmt.Printf("Error with Set call: %s\n", err)
      return err
   }
   defer resp.Body.Close()
   retval, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      fmt.Printf("Error with Set response: %s\n", err)
      return err
   }

   if resp.StatusCode != 200 && resp.StatusCode != 201 {
      fmt.Printf("Error: key = %s (%d) %s\n", key, resp.StatusCode, retval)
      valstr := string(retval)
      return errors.New(valstr)
   }
   
   return nil
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
