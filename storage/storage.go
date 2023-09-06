package storage

import "time"

type Service interface {
  Save(string, time.Time) (string, error)
  Load(string) (string, error)
  LoadInfo(string) (*Item, error)
  Close() error
}

type Item struct {
  Id      uint64 `json:"id" redis:"id"`
  URL     string `json:"url" redis:"url"`
  Expires string `json:"expires" redis:"expires"`
  Visits  int    `json:"visits" redis:"visits"`
}

type redis struct{ pool *redisClient.Pool }

func New(host, port, password string) (storage.Service, error) {
  pool := &redisClient.Pool{
     MaxIdle:     10,
     IdleTimeout: 240 * time.Second,
     Dial: func() (redisClient.Conn, error) {
        return redisClient.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
     },
  }

  return &redis{pool}, nil
}