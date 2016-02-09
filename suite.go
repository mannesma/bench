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
   PerfList map[string]*perf.Event
}

func MakeSuite (config *Config) *Suite {
   s := &Suite {
      Config: config,
      Client: client.MakeClient(config.ClientType,
                                config.ServerHost,
                                config.Debug),
      RandGen: rand.New(rand.NewSource(config.Seed)),
      PerfList: make(map[string]*perf.Event),
   }

   if s.Config.BenchType == "read" {
      fmt.Printf("Info: read test\n")
      s.Benchmark = s.bench_read
      s.Setup = s.setup_read
      s.PerfList["Read"] = perf.MakeEvent("Read")
   }

   return s
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


