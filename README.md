# Onspring API Go SDK

⚠️ Under Construction ⚠️

The Go SDK for the Onspring API is meant to simplify development in Go for Onspring customers who want to build integrations with their Onspring instance.

Note: This is an unofficial SDK for the Onspring API. It was not built in consultation with Onspring Technologies LLC or a member of their development team.

This SDK was developed independently using Onspring's existing C# SDK, the Onspring API's swagger page, and api documentation as the starting point with the intention of making development of integrations done in Javascript with an Onspring instance quicker and more convenient.

## Dependencies

### Go

Requires use of [Go](https://golang.org/dl/) version 1.18 or higher.

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
import "github.com/StevanFreeborn/onspring-api-sdk-go/onspring"

client := onspring.NewClient("your-api-key")
```

The `Client` instance can be further configured by providing optional configuration settings via the `ClientConfig` struct. For example, you can set a custom base URL for the Onspring API if needed:

```go
import "github.com/StevanFreeborn/onspring-api-sdk-go/onspring"

client := onspring.NewClient(
    "your-api-key",
    onspring.WithBaseURL(customURL),
)
```

### Full API Documentation

You may wish to refer to the full [Onspring API documentation](https://software.onspring.com/hubfs/Training/Admin%20Guide%20-%20v2%20API.pdf) when determining which values to pass as parameters to some of the `OnspringClient` methods. There is also a [swagger page](https://api.onspring.com/swagger/index.html) that you can use for making exploratory requests.

## Examples

### Connectivity

#### Verify connectivity

```go
import (
    "fmt"
    "github.com/StevanFreeborn/onspring-api-sdk-go/onspring"
)

client := onspring.NewClient("your-api-key")

err := client.Ping.Get(context.TODO())

if err == nil {
    fmt.Println("Connection successful!")
}
```

