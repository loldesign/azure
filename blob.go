package azure

import (
	"encoding/xml"
	"fmt"
	"github.com/loldesign/azure/core"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
)

var client = &http.Client{}

type Azure struct {
	Account string
	Key     string
}

type Blobs struct {
	XMLName xml.Name `xml:"EnumerationResults"`
	Itens   []Blob   `xml:"Blobs>Blob"`
}

type Blob struct {
	Name     string   `xml:"Name"`
	Property Property `xml:"Properties"`
}

type Property struct {
	LastModified  string `xml:"Last-Modified"`
	Etag          string `xml:"Etag"`
	ContentLength string `xml:"Content-Length"`
	ContentType   string `xml:"Content-Type"`
	BlobType      string `xml:"BlobType"`
	LeaseStatus   string `xml:"LeaseStatus"`
}

func (a Azure) doRequest(azureRequest core.AzureRequest) (*http.Response, error) {
	client, req := a.clientAndRequest(azureRequest)
	return client.Do(req)
}

func (a Azure) clientAndRequest(azureRequest core.AzureRequest) (*http.Client, *http.Request) {
	req := a.prepareRequest(azureRequest)

	return client, req
}

func (a Azure) prepareRequest(azureRequest core.AzureRequest) *http.Request {
	credentials := core.Credentials{
		Account:   a.Account,
		AccessKey: a.Key}

	return core.New(credentials, azureRequest).PrepareRequest()
}

func New(account, accessKey string) Azure {
	return Azure{account, accessKey}
}

func (a Azure) CreateContainer(container string, meta map[string]string) (*http.Response, error) {
	azureRequest := core.AzureRequest{
		Method:      "put",
		Container:   container,
		Resource:    "?restype=container",
		Header:      meta,
		RequestTime: time.Now().UTC()}

	return a.doRequest(azureRequest)
}

func (a Azure) DeleteContainer(container string) (*http.Response, error) {
	azureRequest := core.AzureRequest{
		Method:      "delete",
		Container:   container,
		Resource:    "?restype=container",
		RequestTime: time.Now().UTC()}

	return a.doRequest(azureRequest)
}

func (a Azure) FileUpload(container, name string, body io.Reader) (*http.Response, error) {
	extension := strings.ToLower(path.Ext(name))
	contentType := mime.TypeByExtension(extension)

	azureRequest := core.AzureRequest{
		Method:      "put",
		Container:   container,
		Blob:        name,
		Body:        body,
		Header:      map[string]string{"x-ms-blob-type": "BlockBlob", "Accept-Charset": "UTF-8", "Content-Type": contentType},
		RequestTime: time.Now().UTC()}

	return a.doRequest(azureRequest)
}

func (a Azure) ListBlobs(container string) (Blobs, error) {
	var blobs Blobs

	azureRequest := core.AzureRequest{
		Method:      "get",
		Container:   container,
		Resource:    "?restype=container&comp=list",
		RequestTime: time.Now().UTC()}

	res, err := a.doRequest(azureRequest)

	if err != nil {
		return blobs, err
	}

	decoder := xml.NewDecoder(res.Body)
	decoder.Decode(&blobs)

	return blobs, nil
}

func (a Azure) DeleteBlob(container, name string) (bool, error) {
	azureRequest := core.AzureRequest{
		Method:      "delete",
		Container:   container,
		Blob:        name,
		RequestTime: time.Now().UTC()}

	res, err := a.doRequest(azureRequest)

	if err != nil {
		return false, err
	}

	if res.StatusCode != 202 {
		return false, fmt.Errorf("deleteBlob: %s", res.Status)
	}

	return true, nil
}

func (a Azure) FileDownload(container, name string) (*http.Response, error) {
	azureRequest := core.AzureRequest{
		Method:      "get",
		Container:   container,
		Blob:        name,
		RequestTime: time.Now().UTC()}

	return a.doRequest(azureRequest)
}

func (a Azure) CopyBlob(container, name, source string) (*http.Response, error) {
	azureRequest := core.AzureRequest{
		Method:      "put",
		Container:   container,
		Blob:        name,
		Header:      map[string]string{"x-ms-copy-source": source},
		RequestTime: time.Now().UTC()}

	return a.doRequest(azureRequest)
}
