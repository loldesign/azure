package core

import(
  "fmt"
  "testing"
  "time"
  "net/http"
  "net/http/httptest"
  . "launchpad.net/gocheck"
)

func Test(t *testing.T) {
  TestingT(t)
}

var _ = Suite(&CoreSuite{})

type CoreSuite struct{
  core *Core
}

func (s *CoreSuite) SetUpSuite (c *C) {
  time := time.Date(2013, time.November, 22, 15, 0, 0, 0, time.UTC)
  credentials := Credentials{
    Account: "sampleAccount",
    AccessKey: "secretKey"}

  azureRequest := AzureRequest{
    Method: "put",
    Container: "samplecontainer",
    Resource: "?restype=container",
    RequestTime: time}

  s.core = New(credentials, azureRequest)
}

func (s *CoreSuite) Test_RequestUrl(c *C) {
  expected := "https://sampleAccount.blob.core.windows.net/samplecontainer?restype=container"

  c.Assert(s.core.RequestUrl(), Equals, expected)
}

func (s *CoreSuite) Test_Request(c *C) {
  handle := func(w http.ResponseWriter, r *http.Request) {
    c.Assert(r.URL.Scheme, Equals, "https")
    c.Assert(r.URL.Host, Equals, "sampleAccount.blob.core.windows.net")
    c.Assert(r.URL.Path, Equals, "/samplecontainer")

    //METHOD
    c.Assert(r.Method, Equals, "PUT")
    // HEADER
    c.Assert(r.Header.Get("x-ms-date"), Equals, "Fri, 22 Nov 2013 15:00:00 GMT")
    c.Assert(r.Header.Get("x-ms-version"), Equals, "2009-09-19")
    c.Assert(r.Header.Get("Authorization"), Equals, "SharedKey sampleAccount:5Mb0CpXTPSuX69+4njLQJB9Bf7aoBrsFhamFb7ZHRWs=")
  }

  req := s.core.PrepareRequest()
  w := httptest.NewRecorder()

  handle(w, req)
}

func (s *CoreSuite) Test_RequestWithCustomHeaders(c *C) {
  handle := func(w http.ResponseWriter, r *http.Request) {
    // HEADER
    c.Assert(r.Header.Get("some"), Equals, "header key")
    c.Assert(r.Header.Get("x-ms-blob-type"), Equals, "BlockBlob")
    c.Assert(r.Header.Get("x-ms-date"), Equals, "Fri, 22 Nov 2013 15:00:00 GMT")
    c.Assert(r.Header.Get("x-ms-version"), Equals, "2009-09-19")
    c.Assert(r.Header.Get("Authorization"), Equals, "SharedKey sampleAccount:KcrJEvtAQwa0NJbObaCqBIgXOEcXAMmxcAuMr4+C+E4=")
  }

  s.core.AzureRequest.Header = map[string]string{
    "x-ms-blob-type":"BlockBlob",
    "some": "header key"}

  req := s.core.PrepareRequest()
  w := httptest.NewRecorder()

  handle(w, req)
}

func (s *CoreSuite) Test_CanonicalizedHeaders(c *C) {
  req, err := http.NewRequest("GET", "http://example.com", nil)

  if err != nil {
    c.Error(err)
  }

  req.Header.Add("nothing", "important")
  req.Header.Add("X-Ms-Version", "2009-09-19")
  req.Header.Add("X-Ms-Date", "Fri, 22 Nov 2013 15:00:00 GMT")
  req.Header.Add("X-Ms-Blob-Type", "BlockBlob")
  req.Header.Add("Content-Type", "text/plain; charset=UTF-8")

  s.core.AzureRequest.Request = req

  expected := fmt.Sprintf("x-ms-blob-type:%s\nx-ms-date:%s\nx-ms-version:%s", "BlockBlob", "Fri, 22 Nov 2013 15:00:00 GMT", "2009-09-19")
  c.Assert(s.core.canonicalizedHeaders(), Equals, expected)
}