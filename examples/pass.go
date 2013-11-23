package main

import(
  "fmt"
  "crypto/hmac"
  "crypto/sha512"
  "crypto/sha256"
  "encoding/base64"
  "crypto/md5"
  "crypto/sha1"
  "io"
)

const pass = "123456"

func Md5() {
  md5 := md5.New()
  io.WriteString(md5, pass)
  fmt.Printf("MD5: %x\n", md5.Sum(nil))
}

func Sha1() {
  sha1 := sha1.New()
  io.WriteString(sha1, pass)
  fmt.Printf("SHA1: %x\n", sha1.Sum(nil))
}

func Sha512() {
  sha512 := sha512.New() // NEW
  sha512.Write([]byte(pass))

  fmt.Printf("SHA512: %x\n", sha512.Sum(nil)) //SUM
}

func Sha256() {
  sha256 := sha256.New() // NEW
  sha256.Write([]byte(pass))

  fmt.Printf("sha256: %x\n", sha256.Sum(nil)) //SUM
}

func Sha512Concat() {
  sha512 := hmac.New(sha512.New, []byte("BannersViewer")) // NEW
  sha512.Write([]byte(pass))

  fmt.Printf("sha512: %x\n", sha512.Sum(nil)) //SUM
  // fmt.Printf("sha512: %s\n", sha512.Sum(nil)) //SUM
}

func Azure() string {
  sha256 := hmac.New(sha256.New, []byte("BannersViewer")) // NEW
  sha256.Write([]byte(pass))

  return base64.StdEncoding.EncodeToString(sha256.Sum(nil)) //SUM
}


func main() {
  // Md5() // OK
  // Sha1() // OK
  // Sha256()
  // Sha512()
  // Sha512Concat() // OK!
  fmt.Printf("Azure [sha256]: %s\n", Azure()) // OK!
}