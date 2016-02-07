package main

import (
	"fmt"
   "github.com/mannesma/bench"
	"os"
)

func main() {
   config := bench.MakeConfigFromCmdline()
   t := bench.MakeSuite(config)   

   if t == nil {
      fmt.Printf("Error creating Suite!\n")
      os.Exit(1)
   }

   if t.Config.Setup {
      err := t.Setup()
      if err != nil {
         fmt.Printf("Error with setup: %s\n", err)
      }
   } else {
      t.Run()
      t.ReadPerf.Print(true)
   }
}
