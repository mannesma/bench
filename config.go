package bench

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
   ServerType string
   ServerHost string
   TestType string
   Setup bool
   Iterations int64
   ArrivalRate float64
   Seed int64
   Debug bool
}


func MakeConfigFromCmdline() *Config {
   tc := &Config {}
   flag.StringVar(&tc.ServerType, "server_type", "consul", "Type of server to connect to")   
   flag.StringVar(&tc.TestType, "test_type", "read", "Type of test to run")
   flag.BoolVar(&tc.Setup, "setup", false, "Initialize the servers for test type")
   flag.Int64Var(&tc.Iterations, "iterations", 10, "Number of times to read")
   flag.Float64Var(&tc.ArrivalRate, "arrival_rate", 2, "Number of operations per second")
   flag.Int64Var(&tc.Seed, "seed", time.Now().UnixNano(), "Random number seed (defaults to current nanoseconds)")
   flag.BoolVar(&tc.Debug, "debug", false, "Enable verbose output")
   flag.StringVar(&tc.ServerHost, "server_host", "", "Override of server host:port")

   flag.Parse()

   if tc.ServerType != "consul" && 
      tc.ServerType != "etcd" && 
      tc.ServerType != "zookeeper" {
      fmt.Printf("Error: invalid server_type '%s'\n", tc.ServerType)
   }
   fmt.Printf("server_type = %s\n", tc.ServerType)
   fmt.Printf("server_host = %s\n", tc.ServerHost)
   fmt.Printf("setup = %t\n", tc.Setup)
   fmt.Printf("test_type = %s\n", tc.TestType)
   fmt.Printf("iterations = %d\n", tc.Iterations)
   fmt.Printf("arrival_rate = %f\n", tc.ArrivalRate)
   fmt.Printf("seed = %d\n", tc.Seed)
   fmt.Printf("debug = %t\n", tc.Debug)

   return tc
}
