package core

import(
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
  s.core = New("sampleAccount", "secretKey", "PUT", "samplecontainer", "restype=container", time)
}

func (s *CoreSuite) Test_RequestUrl(c *C) {
  expected := "http://sampleAccount.blob.core.windows.net/samplecontainer?restype=container"

  c.Assert(s.core.RequestUrl(), Equals, expected)
}

func (s *CoreSuite) Test_Request(c *C) {
  handle := func(w http.ResponseWriter, r *http.Request) {
    c.Assert(r.Header.Get("x-ms-date"), Equals, "Fri, 22 Nov 2013 15:00:00 GMT")
    c.Assert(r.Header.Get("x-ms-version"), Equals, "2009-09-19")
    c.Assert(r.Header.Get("Authorization"), Equals, "SharedKey sampleAccount:5Mb0CpXTPSuX69+4njLQJB9Bf7aoBrsFhamFb7ZHRWs=")
  }

  req := s.core.Request()
  w := httptest.NewRecorder()

  handle(w, req)
}