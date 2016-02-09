package main

import (
	"fmt"
   "github.com/mannesma/bench"
	"os"
)

func main() {
   config := bench.MakeConfigFromCmdline()
   s := bench.MakeSuite(config)   

   if s == nil {
      fmt.Printf("Error creating Suite!\n")
      os.Exit(1)
   }

   if s.Config.Setup {
      err := s.Setup()
      if err != nil {
         fmt.Printf("Error with setup: %s\n", err)
      }
   } else {
      s.Run()
      for k, v := range s.PerfList {
         fmt.Printf("bench = %s\n", k)
         v.Print(true)
      }
   }
}
