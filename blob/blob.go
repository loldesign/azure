package blob

import (
  "github.com/loldesign/azure/core"
  "net/http"
  "time"
  "log"
)

type Azure struct {
  Account string
  AccessKey string
}

func (azure Azure) prepareRequest(method, container, resource string) *http.Request {
  core := core.Core{
      AccessKey: azure.AccessKey,
      Account: azure.Account,
      Method: method,
      RequestTime: time.Now().UTC(),
      Container: container,
      Resource: resource}

  return core.PrepareRequest()
}

func (azure Azure) CreateContainer(container string) *http.Response {
  client := &http.Client{}
  req := azure.prepareRequest("put", container, "?restype=container")

  res, err := client.Do(req)

  if err != nil {
    log.Fatal(err)
  }

  return res
}

func (azure Azure) DeleteContainer(container string) *http.Response {
  client := &http.Client{}
  req := azure.prepareRequest("delete", container, "?restype=container")

  res, err := client.Do(req)

  if err != nil {
    log.Fatal(err)
  }

  return res
}
