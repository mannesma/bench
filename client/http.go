package client

import (
   "bytes"
   "errors"
   "fmt"
   "io/ioutil"
   "net/http"
)

type HttpClient struct {
   base_url string
   client *http.Client
}

func MakeHttpClient(base_url string) *HttpClient {

   client := &HttpClient {
      base_url: base_url,
      client: &http.Client{},
   }

   return client
}

func (h *HttpClient) Get(key string) ([]byte, error) {
   resp, err := h.client.Get(h.base_url + key)
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

func (h *HttpClient) Put(key string, value []byte) error {
   br := bytes.NewReader(value)
   req, err := http.NewRequest("PUT", h.base_url + key, br)
   if err != nil {
      fmt.Printf("Error with Set request: %s\n", err)
      return err
   }
   resp, err := h.client.Do(req)
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
