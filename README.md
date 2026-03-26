# Onspring API Go SDK

[![Pull Request](https://github.com/StevanFreeborn/onspring-api-sdk-go/actions/workflows/pull_request.yml/badge.svg)](https://github.com/StevanFreeborn/onspring-api-sdk-go/actions/workflows/pull_request.yml)
[![Version](https://github.com/StevanFreeborn/onspring-api-sdk-go/actions/workflows/version.yml/badge.svg)](https://github.com/StevanFreeborn/onspring-api-sdk-go/actions/workflows/version.yml)
[![Publish](https://github.com/StevanFreeborn/onspring-api-sdk-go/actions/workflows/publish.yml/badge.svg)](https://github.com/StevanFreeborn/onspring-api-sdk-go/actions/workflows/publish.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/StevanFreeborn/onspring-api-sdk-go.svg)](https://pkg.go.dev/github.com/StevanFreeborn/onspring-api-sdk-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.md)

The Go SDK for the Onspring API is meant to simplify development in Go for Onspring customers who want to build integrations with their Onspring instance.

Note: This is an unofficial SDK for the Onspring API. It was not built in consultation with Onspring Technologies LLC or a member of their development team.

This SDK was developed independently using Onspring's existing C# SDK, the Onspring API's swagger page, and api documentation as the starting point with the intention of making development of integrations done in Go with an Onspring instance quicker and more convenient.

## Dependencies

### Go

Requires use of [Go](https://golang.org/dl/) version 1.25 or higher.

## Installation

To install the Onspring API Go SDK, use the following command:

```pwsh
go get github.com/StevanFreeborn/onspring-api-sdk-go
```

## API Key

In order to successfully interact with the Onspring Api you will need an API key. API keys are obtained by an Onspring user with permissions to at least **Read** API Keys for your instance via the following steps:

1. Login to the Onspring instance.
2. Navigate to **Administration** > **Security** > **API Keys**
3. On the list page, add a new API Key - this will require **Create** permissions - or click an existing API key to view its details.
4. Click on the **Developer Information** tab.
5. Copy the **X-ApiKey Header** value from this tab.

**Important:**

- An API Key must have a status of `Enabled` in order to make authorized requests.
- Each API Key must have an assigned Role. This role controls the permissions for requests made. If the API Key used does not have sufficient permissions the requests made won't be successful.

### 🔒 Permission Considerations

You can think of any API Key as another user in your Onspring instance and therefore it is subject to all the same permission considerations as any other user when it comes to its ability to access data in your instance. The API Key you use needs to have all the correct permissions within your instance to access the data requested. Things to think about in this context are `role security`, `content security`, and `field security`.

## Start Coding

### `Client`

The most common way to use the SDK is to create a `Client` instance and call its methods to interact with the Onspring API. Here is an example of how to create a `Client` instance. You will need to provide your API key when creating the client. It is best practice to store your API key securely and not hard-code it in your source code.

```go
import "github.com/StevanFreeborn/onspring-api-sdk-go"

client := onspring.NewClient("your-api-key")
```

The `Client` instance can be further configured by providing optional functional options. For example, you can set a custom base URL for the Onspring API if needed:

```go
import "github.com/StevanFreeborn/onspring-api-sdk-go"

client := onspring.NewClient(
  "your-api-key",
  onspring.WithBaseURL(customURL),
)
```

### Full API Documentation

You may wish to refer to the full [Onspring API documentation](https://software.onspring.com/hubfs/Training/Admin%20Guide%20-%20v2%20API.pdf) when determining which values to pass as parameters to some of the `Client` methods. There is also a [swagger page](https://api.onspring.com/swagger/index.html) that you can use for making exploratory requests.

## Examples

### Connectivity

#### Verify connectivity

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

err := client.Ping.Get(context.TODO())

if err == nil {
  fmt.Println("Connection successful!")
}
```

### Apps

#### Get App by Id

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

app, err := client.Apps.Get(context.TODO(), 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved App: %+v\n", app)
}
```

#### Get Apps by Page

##### Retrieve a single page

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

page, err := client.Apps.List(context.TODO())

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Page %d of Apps: %+v\n", page.PageNumber, page.Items)
}
```

##### Retrieve all pages

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

for app, err := range client.Apps.ListAll(context.TODO()) {
  if err != nil {
    fmt.Printf("Error during iteration: %v\n", err)
    break
  }
  fmt.Printf("Retrieved App: %+v\n", app)
}
```

#### Get Apps by Batch

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

batch, err := client.Apps.GetMany(context.TODO(), []int{ 1, 2 })

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved App Batch (Count: %d): %+v\n", batch.Count, batch.Items)
}
```

### Fields

#### Get Field by Id

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

field, err := client.Fields.Get(context.TODO(), 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Field: %+v\n", field)
}
```

#### Get Fields by App

##### Retrieve a single page

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

page, err := client.Fields.List(context.TODO(), 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Page %d of Fields: %+v\n", page.PageNumber, page.Items)
}
```

##### Retrieve all pages

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

for field, err := range client.Fields.ListAll(context.TODO(), 1) {
  if err != nil {
    fmt.Printf("Error during iteration: %v\n", err)
    break
  }
  fmt.Printf("Retrieved Field: %+v\n", field)
}
```

#### Get Fields by Batch

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

batch, err := client.Fields.GetMany(context.TODO(), []int{ 1, 2 })

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Field Batch (Count: %d): %+v\n", batch.Count, batch.Items)
}
```

### Lists

#### Save List Item

Create a new list item or update an existing one. To create a new item, omit the `Id` field. To update an existing item, provide the `Id` of the item to update.

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

numericValue := 1.0
color := "#ff0000"

item := onspring.SaveListItemRequest{
  Name:         "New List Value",
  NumericValue: &numericValue,
  Color:        &color,
}

response, err := client.Lists.Save(context.TODO(), 1, item)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Saved List Item Id: %s\n", response.Id)
}
```

#### Delete List Item

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

err := client.Lists.Delete(context.TODO(), 1, "d4a3c2b1-e5f6-7890-abcd-ef1234567890")

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Println("List item deleted successfully!")
}
```

### Reports

#### Get Report by Id

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

report, err := client.Reports.Get(context.TODO(), 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Report Columns: %v\n", report.Columns)
  fmt.Printf("Retrieved Report Rows: %d\n", len(report.Rows))
}
```

You can also specify the data format and data type for the report:

```go
report, err := client.Reports.Get(
  context.TODO(),
  1,
  onspring.WithDataFormat("Formatted"),
  onspring.WithDataType("ChartData"),
)
```

#### Get Reports by App

##### Retrieve a single page

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

page, err := client.Reports.List(context.TODO(), 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Page %d of Reports: %+v\n", page.PageNumber, page.Items)
}
```

##### Retrieve all pages

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

for report, err := range client.Reports.ListAll(context.TODO(), 1) {
  if err != nil {
    fmt.Printf("Error during iteration: %v\n", err)
    break
  }
  fmt.Printf("Retrieved Report: %+v\n", report)
}
```

### Files

#### Get File Info

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

fileInfo, err := client.Files.GetInfo(context.TODO(), 1, 2, 3)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("File Name: %s\n", fileInfo.Name)
  fmt.Printf("Content Type: %s\n", fileInfo.ContentType)
}
```

#### Get File Content

```go
import (
  "context"
  "fmt"
  "os"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

fileContent, err := client.Files.GetContent(context.TODO(), 1, 2, 3)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("File Name: %s\n", fileContent.FileName)
  fmt.Printf("Content Type: %s\n", fileContent.ContentType)
  os.WriteFile(fileContent.FileName, fileContent.Data, 0644)
}
```

#### Delete File

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

err := client.Files.Delete(context.TODO(), 1, 2, 3)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Println("File deleted successfully!")
}
```

#### Save File

```go
import (
  "context"
  "fmt"
  "os"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

file, _ := os.Open("document.pdf")
defer file.Close()

saveReq := onspring.SaveFileRequest{
  RecordId:     1,
  FieldId:      2,
  FileName:     "document.pdf",
  FileContents: file,
  Notes:        "Uploaded via API",
  ModifiedDate: "2024-01-01T00:00:00Z",
}

response, err := client.Files.Save(context.TODO(), saveReq)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Saved File Id: %d\n", response.Id)
}
```

### Records

#### Get Record by Id

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

record, err := client.Records.Get(context.TODO(), 1, 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Record: %+v\n", record)
}
```

You can also specify which fields to include and the data format:

```go
record, err := client.Records.Get(
  context.TODO(), 1, 1,
  onspring.WithFieldIds([]int{1, 2, 3}),
  onspring.WithRecordDataFormat("Formatted"),
)
```

#### Get Records by App

##### Retrieve a single page

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

page, err := client.Records.List(context.TODO(), 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Page %d of Records: %+v\n", page.PageNumber, page.Items)
}
```

You can combine paging, field, and data format options:

```go
page, err := client.Records.List(
  context.TODO(), 1,
  onspring.WithFieldIds([]int{1, 2}),
  onspring.WithRecordDataFormat("Formatted"),
  onspring.WithPaging(onspring.ForPageNumber(2), onspring.WithPageSize(10)),
)
```

##### Retrieve all pages

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

for record, err := range client.Records.ListAll(context.TODO(), 1) {
  if err != nil {
    fmt.Printf("Error during iteration: %v\n", err)
    break
  }
  fmt.Printf("Retrieved Record: %+v\n", record)
}
```

#### Get Records by Batch

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

batch, err := client.Records.GetMany(context.TODO(), onspring.GetManyRecordsRequest{
  AppId:     1,
  RecordIds: []int{1, 2},
  FieldIds:  []int{1, 2},
})

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Record Batch (Count: %d): %+v\n", batch.Count, batch.Items)
}
```

#### Query Records

##### Retrieve a single page

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

page, err := client.Records.Query(context.TODO(), onspring.QueryRecordsRequest{
  AppId:  1,
  Filter: "field eq 'value'",
})

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Retrieved Page %d of Records: %+v\n", page.PageNumber, page.Items)
}
```

##### Retrieve all pages

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

for record, err := range client.Records.QueryAll(context.TODO(), onspring.QueryRecordsRequest{
  AppId:  1,
  Filter: "field eq 'value'",
}) {
  if err != nil {
    fmt.Printf("Error during iteration: %v\n", err)
    break
  }
  fmt.Printf("Retrieved Record: %+v\n", record)
}
```

#### Save Record

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

response, err := client.Records.Save(context.TODO(), onspring.SaveRecordRequest{
  AppId:  1,
  Fields: map[string]any{"1": "value", "2": 42},
})

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Printf("Saved Record Id: %d\n", response.Id)
}
```

#### Delete Record

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

err := client.Records.Delete(context.TODO(), 1, 1)

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Println("Record deleted successfully!")
}
```

#### Delete Records by Batch

```go
import (
  "context"
  "fmt"
  "github.com/StevanFreeborn/onspring-api-sdk-go"
)

client := onspring.NewClient("your-api-key")

err := client.Records.DeleteMany(context.TODO(), onspring.DeleteManyRecordsRequest{
  AppId:     1,
  RecordIds: []int{1, 2, 3},
})

if err != nil {
  fmt.Printf("Error: %v\n", err)
} else {
  fmt.Println("Records deleted successfully!")
}

