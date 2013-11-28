package azure

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

type azure struct {
  account string
  key string
}

func (a azure) doRequest(azureRequest core.AzureRequest) (*http.Response, error) {
  client, req := a.clientAndRequest(azureRequest)
  return client.Do(req)
}

func (a azure) clientAndRequest(azureRequest core.AzureRequest) (*http.Client, *http.Request) {
  req := a.prepareRequest(azureRequest)

  return client, req
}

func (a azure) prepareRequest(azureRequest core.AzureRequest) *http.Request {
  credentials := core.Credentials{
    Account: a.account,
    AccessKey: a.key}

  return core.New(credentials, azureRequest).PrepareRequest()
}

func prepareMetadata(keys map[string]string) map[string]string {
  header := make(map[string]string)

  for k, v := range keys {
    key := fmt.Sprintf("x-ms-meta-%s", k)
    header[key] = v
  }

  return header
}

func New(account, accessKey string) azure {
  return azure{account, accessKey}
}

func (a azure) CreateContainer(container string, meta map[string]string) (*http.Response, error) {
  azureRequest := core.AzureRequest{
    Method: "put",
    Container: container,
    Resource: "?restype=container",
    Header: prepareMetadata(meta),
    RequestTime: time.Now().UTC()}

  return a.doRequest(azureRequest)
}

func (a azure) DeleteContainer(container string) (*http.Response, error) {
  azureRequest := core.AzureRequest{
    Method: "delete",
    Container: container,
    Resource: "?restype=container",
    RequestTime: time.Now().UTC()}

  return a.doRequest(azureRequest)
}

func (a azure) FileUpload(container, name string, file *os.File) (*http.Response, error) {
  extension := strings.ToLower(path.Ext(file.Name()))
  contentType := mime.TypeByExtension(extension)

  azureRequest := core.AzureRequest{
    Method: "put",
    Container: fmt.Sprintf("%s/%s", container, name),
    Body: file,
    Header: map[string]string{"x-ms-blob-type": "BlockBlob", "Accept-Charset": "UTF-8", "Content-Type": contentType},
    RequestTime: time.Now().UTC()}

  return a.doRequest(azureRequest)
}