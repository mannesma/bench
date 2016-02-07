package bench

import (
   "fmt"
   "github.com/mannesma/bench/client"
   "github.com/mannesma/bench/perf"
   "math/rand"
   "time"
)

type Suite struct {
   Config *Config
   Client client.Client
   RandGen *rand.Rand
   Benchmark func() error
   Setup func() error
   ReadPerf *perf.Event
}

func MakeSuite (config *Config) *Suite {
   s := &Suite {
      Config: config,
      Client: client.MakeClient(config.ClientType, config.ServerHost),
      RandGen: rand.New(rand.NewSource(config.Seed)),
      ReadPerf: perf.MakeEvent("Get"),
   }

   if s.Config.BenchType == "read" {
      fmt.Printf("Info: read test\n")
      s.Benchmark = s.bench_read
      s.Setup = s.setup_read
   }

   return s
}

func (s *Suite) setup_read() error {
   for _, k1 := range KeyNames {
      key := fmt.Sprintf("/%s", k1)
      err := s.Client.Set(key, []byte(fmt.Sprintf("%d", s.RandGen.Intn(1024))))
      if err != nil {
         return err
      }
      for _, k2 := range KeyNames {
         key := fmt.Sprintf("/%s/%s", k1, k2)
         err := s.Client.Set(key, []byte(fmt.Sprintf("%d", s.RandGen.Intn(1024))))
         if err != nil {
            return err
         }
      }
   }

   return nil
}

func (s *Suite) bench_read() error {
   key1 := KeyNames[s.RandGen.Intn(len(KeyNames))]
   key2 := KeyNames[s.RandGen.Intn(len(KeyNames))]
   key := fmt.Sprintf("/%s/%s", key1, key2)
   s.ReadPerf.Start()
   value, err := s.Client.Get(key)
   s.ReadPerf.Stop()
   if err == nil {
      if s.Config.Debug { fmt.Printf("value = %s\n", value) }
      return nil
   } else {
      return err
   }
}

func (s *Suite) Run() {
   var i int64
   var sleepval time.Duration

   for i = 0; i < s.Config.Iterations; i++ {
      sleepval = s.calc_sleep_time()
      time.Sleep(sleepval)
      s.Benchmark()
   }
}

func (s *Suite) calc_sleep_time() time.Duration {
   var sleepval float64
   sleepval = s.RandGen.ExpFloat64() / s.Config.ArrivalRate
   if s.Config.Debug { fmt.Printf("calc sleepval = %f\n", sleepval) }
   duration := time.Duration(int64(sleepval * float64(time.Second)))
   if s.Config.Debug { fmt.Printf("duration_in = %T %v\n", duration, duration) }
   return duration
}


