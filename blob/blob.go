package blob

import (
  "github.com/loldesign/azure/core"
  "net/http"
  "fmt"
  "time"
  "os"
  "mime"
  "strings"
  "path"
)

var client = &http.Client{}

type Azure struct {
  Account string
  AccessKey string
}

func (azure Azure) doRequest(azureRequest core.AzureRequest) (*http.Response, error) {
  client, req := azure.clientAndRequest(azureRequest)
  return client.Do(req)
}

func (azure Azure) clientAndRequest(azureRequest core.AzureRequest) (*http.Client, *http.Request) {
  req := azure.prepareRequest(azureRequest)

  return client, req
}

func (azure Azure) prepareRequest(azureRequest core.AzureRequest) *http.Request {
  credentials := core.Credentials{
    Account: azure.Account,
    AccessKey: azure.AccessKey}

  return core.New(credentials, azureRequest).PrepareRequest()
}

func (azure Azure) CreateContainer(container string) (*http.Response, error) {
  azureRequest := core.AzureRequest{
    Method: "put",
    Container: container,
    Resource: "?restype=container",
    RequestTime: time.Now().UTC()}

  return azure.doRequest(azureRequest)
}

func (azure Azure) DeleteContainer(container string) (*http.Response, error) {
  azureRequest := core.AzureRequest{
    Method: "delete",
    Container: container,
    Resource: "?restype=container",
    RequestTime: time.Now().UTC()}

  return azure.doRequest(azureRequest)
}

func (azure Azure) FileUpload(container, name string, file *os.File) (*http.Response, error) {
  extension := strings.ToLower(path.Ext(file.Name()))
  contentType := mime.TypeByExtension(extension)

  azureRequest := core.AzureRequest{
    Method: "put",
    Container: fmt.Sprintf("%s/%s", container, name),
    Body: file,
    Header: map[string]string{"x-ms-blob-type": "BlockBlob", "Accept-Charset": "UTF-8", "Content-Type": contentType},
    RequestTime: time.Now().UTC()}

  return azure.doRequest(azureRequest)
}