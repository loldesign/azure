package core

import(
  "time"
  "net/url"
  "net/http"
  "fmt"
  "log"
  "bytes"
  "strings"
  "encoding/base64"
  "crypto/hmac"
  "crypto/sha256"
  "sort"
)

const ms_date_layout = "Mon, 2 Jan 2006 15:04:05 GMT"
const version = "2009-09-19"

type Core struct {
  Account string
  AccessKey string
  Method string
  Container string
  Resource string
  RequestTime time.Time
}

func New(account, accessKey, method, container, resource string, requestTime time.Time) *Core {
  return &Core{
    Account: account,
    AccessKey: accessKey,
    Method: method,
    Container: container,
    Resource: resource,
    RequestTime: requestTime}
}

func (core Core) WebService() string {
  return fmt.Sprintf("http://%s.blob.core.windows.net/", core.Account)
}

func (core Core) RequestUrl() string {
  return fmt.Sprintf("%s%s?%s", core.WebService(), core.Container, core.Resource)
}

func (core Core) Request() *http.Request {
  req, err := http.NewRequest(strings.ToUpper(core.Method), core.RequestUrl(), nil)

  if err != nil {
    log.Fatal(err)
  }

  core.addHeaderInformations(req)

  return req
}

func (core Core) addHeaderInformations(req *http.Request) {
  req.Header.Add("x-ms-date", core.FormattedRequestTime())
  req.Header.Add("x-ms-version", version)
  req.Header.Add("Authorization", core.authorizationHeader())
}

func (core Core) FormattedRequestTime() string {
  return core.RequestTime.Format(ms_date_layout)
}

func (core Core) canonicalizedHeaders() string {
  return fmt.Sprintf("x-ms-date:%s\nx-ms-version:%s", core.FormattedRequestTime(), version)
}

func (core Core) canonicalizedResource() string {
  var buffer bytes.Buffer

  u, err := url.Parse(core.RequestUrl())

  if err != nil {
    log.Fatal(err)
  }

  buffer.WriteString(fmt.Sprintf("/%s/%s", core.Account, core.Container))
  queries := u.Query()

  for key, values := range queries {
    sort.Strings(values)
    buffer.WriteString(fmt.Sprintf("\n%s:%s", key, strings.Join(values, ",")))
  }

  return buffer.String()
}

/*
params:
 HTTP Verb
 Content-Encoding
 Content-Language
 Content-Length
 Content-MD5
 Content-Type
 Date
 If-Modified-Since
 If-Match
 If-None-Match
 If-Unmodified-Since
 Range
*/
func (core Core) signature() string {
  signature := fmt.Sprintf("%s\n\n\n%s\n\n\n\n\n\n\n\n\n%s\n%s",
    core.Method,
    "0",
    core.canonicalizedHeaders(),
    core.canonicalizedResource())

  decodedKey, _ := base64.StdEncoding.DecodeString(core.AccessKey)

  sha256 := hmac.New(sha256.New, []byte(decodedKey))
  sha256.Write([]byte(signature))

  return base64.StdEncoding.EncodeToString(sha256.Sum(nil))
}

func (core Core) authorizationHeader() string {
  return fmt.Sprintf("SharedKey %s:%s", core.Account, core.signature())
}
