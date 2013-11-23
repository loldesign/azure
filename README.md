# Azure

A golang API to communicate with the Azure Storage.
For while, only for manager blobs and containers (create, destroy and so on).

## Installation

```go get github.com/loldesign/azure```

### Usage

#### Creating a container

```go
import(
  "fmt"
  "github.com/loldesign/azure"
)

blob := blob.Azure{Account: "accountName", AccessKey: "secret"}
res := blob.CreateContainer("mycontainer")

fmt.Printf("status -> %s", res.Status)
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am "Added some feature"`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request