# Azure [![Build Status](https://travis-ci.org/loldesign/azure.png)](https://travis-ci.org/loldesign/azure)

A golang API to communicate with the Azure Storage.
For while, only for manager blobs and containers (create, destroy and so on).

## Installation

```go get github.com/loldesign/azure```

## Usage

### Creating a container

```go
package main

import(
  "fmt"
  "github.com/loldesign/azure"
)

func main() {
  blob := azure.New("accountName", "secret")
  res, err := blob.CreateContainer("mycontainer")

  if err != nil {
    fmt.Println(err)
  }

  fmt.Printf("status -> %s", res.Status)
}
```

### Uploading a file to container

```go
package main

import(
  "fmt"
  "github.com/loldesign/azure"
)

func main() {
  blob := azure.New("accountName", "secret")

  file, err := os.Open("path/of/myfile.txt")

  if err != nil {
    fmt.Println(err)
  }

  res, err := blob.FileUpload("mycontainer", "file_name.txt", file)

  if err != nil {
    fmt.Println(err)
  }

  fmt.Printf("status -> %s", res.Status)
}
```

### Listing container's blobs

```go
package main

import(
  "fmt"
  "github.com/loldesign/azure"
)

func main() {
  blob := azure.New("accountName", "secret")
  blobs, err := blob.ListBlobs("mycontainer")

  if err != nil {
    fmt.Println(err)
  }

  for _, file := range blobs.Itens {
    fmt.Printf("blob -> %+v", file)
  }
}
```

### Deleting a blob

```go
package main

import(
  "fmt"
  "github.com/loldesign/azure"
)

func main() {
  blob := azure.New("accountName", "secret")
  ok, err := blob.DeleteBlob("mycontainer", "my_file.png")

  if err != nil {
    fmt.Println(err)
  }

  fmt.Printf("deleted? -> %t", ok)
}
```

### Deleting a container

```go
package main

import(
  "fmt"
  "github.com/loldesign/azure"
)

func main() {
  blob := azure.New("accountName", "secret")
  res, err := blob.DeleteContainer("mycontainer")

  if err != nil {
    fmt.Println(err)
  }

  fmt.Printf("status -> %s", res.Status)
}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am "Added some feature"`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
