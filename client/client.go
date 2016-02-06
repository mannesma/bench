package client

type Client interface {
   Get(key string) ([]byte, error)
   Set(key string, value []byte) error
}
