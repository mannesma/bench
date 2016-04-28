package bench

import (
   "fmt"
)

func (s *Suite) setup_read() error {
   var err error
   for _, k1 := range KeyNames {
      key := fmt.Sprintf("/%s", k1)
      if s.Config.ClientType == "zookeeper" {
         err = s.Client.CreateDir(key)
      }
      if err != nil {
         return err
      }
      for _, k2 := range KeyNames {
         key := fmt.Sprintf("/%s/%s", k1, k2)
         err = s.Client.Set(key, []byte(fmt.Sprintf("%d", s.RandGen.Intn(1024))))
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
   s.PerfList["Read"].Start()
   value, err := s.Client.Get(key)
   s.PerfList["Read"].Stop()
   if err == nil {
      if s.Config.Debug { fmt.Printf("value = %s\n", value) }
      return nil
   } else {
      return err
   }
}
