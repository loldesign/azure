package blob

import (
  "github.com/loldesign/azure/core"
  "fmt"
  "net/http"
  "time"
  "log"
)

type Azure struct {
  Account string
  AccessKey string
}

func (azure Azure) CreateContainer(container string) *http.Response {
  core := core.Core{
    AccessKey: azure.AccessKey,
    Account: azure.Account,
    Method: "PUT", 
    RequestTime: time.Now().UTC(),
    Container: container, 
    Resource: "restype=container"}

  client := &http.Client{}
  req := core.Request()

  res, err := client.Do(req)

  if err != nil {
    log.Fatal(err)
  }

  return res
}
