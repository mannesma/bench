package bench

import (
   "fmt"
   "github.com/mannesma/bench/perf"
   "math/rand"
)

type Suite struct {
   Config *Config
   Client Client
   RandGen *rand.Rand
   TestFunc func() error
   SetupTest func() error
   ReadPerf *perf.Event
}

func MakeSuite (config *Config) *Suite {
   ts := &Suite {
      Config: config,
      RandGen: rand.New(rand.NewSource(config.Seed)),
      ReadPerf: perf.MakeEvent("Get"),
   }

   if config.ServerType == "consul" || config.ServerType == "etcd" {
      ts.Client = MakeTestHttpClient(config)
   } else {  // "zookeeper"
      client, err := MakeTestZookeeperClient(config)
      if err != nil {
         return nil
      }
      ts.Client = client
   }

   if ts.Config.TestType == "read" {
      fmt.Printf("Info: read test\n")
      ts.TestFunc = ts.read_test
      ts.SetupTest = ts.read_setup
   }

   return ts
}

func (s *Suite) read_setup() error {
   for _, k1 := range key_names {
      key := fmt.Sprintf("/%s", k1)
      err := s.Client.Set(key, []byte(fmt.Sprintf("%d", s.RandGen.Intn(1024))))
      if err != nil {
         return err
      }
      for _, k2 := range key_names {
         key := fmt.Sprintf("/%s/%s", k1, k2)
         err := s.Client.Set(key, []byte(fmt.Sprintf("%d", s.RandGen.Intn(1024))))
         if err != nil {
            return err
         }
      }
   }

   return nil
}

func (s *Suite) read_test() error {
   key1 := key_names[s.RandGen.Intn(len(key_names))]
   key2 := key_names[s.RandGen.Intn(len(key_names))]
   key := fmt.Sprintf("/%s/%s", key1, key2)
   s.ReadPerf.Start()
   value, err := s.Client.Get(key)
   s.ReadPerf.Stop()
   if err == nil {
      if ts.Config.Debug { fmt.Printf("value = %s\n", value) }
      return nil
   } else {
      return err
   }
}


func (s *Suite) run_test() {
   var i int64
   var sleepval time.Duration

   for i = 0; i < s.Config.Iterations; i++ {
      sleepval = s.calc_sleep_time()
      time.Sleep(sleepval)
      s.TestFunc()
   }
}

