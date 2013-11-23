package main

import(
  "fmt"
  "time"
  "net/http"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/base64"
)

const azureUrl    = "http://bannersviewers.blob.core.windows.net/"
const dateLayout  = "Mon, 2 Jan 2006 15:04:05 GMT"
const version     = "2009-09-19"
const storageKey  = "zkwQX4LdV6cdWrJMholGuq9MIShlAUjilz+QK+NvtMAv33PCqgOX8pFRsvYws6+L/dqfgXS+mqzjveaJYqYYqQ=="
const accountName = "bannersviewers"
const clientFolder = "marcosinger"

// OK
func canonicalizedHeaders(formattedDate string) string {
  return fmt.Sprintf("x-ms-date:%s\nx-ms-version:%s", formattedDate, version)
}

func canonicalizedResource() string {
  return fmt.Sprintf("/%s/%s\nrestype:container", accountName, clientFolder)
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
func signature(formattedDate string) string {
  signature := fmt.Sprintf("%s\n\n\n%s\n\n\n\n\n\n\n\n\n%s\n%s",
    "PUT",
    "0",
    canonicalizedHeaders(formattedDate),
    canonicalizedResource())

  // EXAMPLE
  // signature := "PUT\n\n\n0\n\n\n\n\n\n\n\n\nx-ms-date:Thu, 21 Nov 2013 16:40:02 GMT\nx-ms-version:2009-09-19\n/bannersviewers/samplecontainerNovo\nrestype:container"
  decodedKey, _ := base64.StdEncoding.DecodeString(storageKey)

  sha256 := hmac.New(sha256.New, []byte(decodedKey)) // NEW
  sha256.Write([]byte(signature))

  return base64.StdEncoding.EncodeToString(sha256.Sum(nil)) //SUM
}

func AuthorizationHeader(formattedDate string) string {
  return fmt.Sprintf("SharedKey %s:%s", accountName, signature(formattedDate))
}

func AddHeaderInformations(req *http.Request, formattedDate string) {
  req.Header.Add("x-ms-date", formattedDate)
  req.Header.Add("x-ms-version", version)
  req.Header.Add("Authorization", AuthorizationHeader(formattedDate))
}

func callAzure(req *http.Request) *http.Response {
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    fmt.Printf("Occurred an error: %s", err.Error())
  }

  return res
}

func main() {
  restype := "container"
  create_container_url := fmt.Sprintf("%s%s?restype=%s", azureUrl, clientFolder, restype)
  
  t := time.Now().UTC()
  date := t.Format(dateLayout)

  req, _ := http.NewRequest("PUT", create_container_url, nil)

  AddHeaderInformations(req, date)
  fmt.Printf("request -> %+v", req)
  fmt.Println(callAzure(req))
  // signature(date)
  // fmt.Println(AuthorizationHeader(date))
}
