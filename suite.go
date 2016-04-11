package bench

import (
	"encoding/json"
   "fmt"
   "github.com/mannesma/bench/client"
   "github.com/mannesma/bench/perf"
   "math/rand"
	"os"
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
		if s.Config.Debug {
      	fmt.Printf("Info: read test\n")
		}
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

func (s *Suite) PrintResults() {
	var result struct {
		Config *Config          `json:Config`
		PerfList []*perf.Report `json:PerfList`
	}

	result.Config = s.Config
	// result.PerfList = make([]*perf.Report)
	for k, v := range s.PerfList {
		r := perf.MakeReport(v)
		if s.Config.Debug {
 			fmt.Printf("bench = %s\n", k)
 			r.Print()
 		} else {
			result.PerfList = append(result.PerfList, r)
		}
	}

	if ! s.Config.Debug {
		str, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("Error: Failed to convert result to json: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", str)
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


