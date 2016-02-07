package bench

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
   ClientType string
   ServerHost string
   BenchType string
   Setup bool
   Iterations int64
   ArrivalRate float64
   Seed int64
   Debug bool
}


func MakeConfigFromCmdline() *Config {
   tc := &Config {}
   flag.StringVar(&tc.ClientType, "client_type", "consul", "Type of client to connect with")   
   flag.StringVar(&tc.BenchType, "bench_type", "read", "Type of test to run")
   flag.BoolVar(&tc.Setup, "setup", false, "Initialize the servers for test type")
   flag.Int64Var(&tc.Iterations, "iterations", 10, "Number of times to read")
   flag.Float64Var(&tc.ArrivalRate, "arrival_rate", 2, "Number of operations per second")
   flag.Int64Var(&tc.Seed, "seed", time.Now().UnixNano(), "Random number seed (defaults to current nanoseconds)")
   flag.BoolVar(&tc.Debug, "debug", false, "Enable verbose output")
   flag.StringVar(&tc.ServerHost, "server_host", "", "Override of server host:port")

   flag.Parse()

   if tc.ClientType != "consul" && 
      tc.ClientType != "etcd" && 
      tc.ClientType != "zookeeper" {
      fmt.Printf("Error: invalid client_type '%s'\n", tc.ClientType)
   }
   fmt.Printf("server_type = %s\n", tc.ClientType)
   fmt.Printf("server_host = %s\n", tc.ServerHost)
   fmt.Printf("setup = %t\n", tc.Setup)
   fmt.Printf("bench_type = %s\n", tc.BenchType)
   fmt.Printf("iterations = %d\n", tc.Iterations)
   fmt.Printf("arrival_rate = %f\n", tc.ArrivalRate)
   fmt.Printf("seed = %d\n", tc.Seed)
   fmt.Printf("debug = %t\n", tc.Debug)

   return tc
}
