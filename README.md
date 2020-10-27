# WebDriver Client (wdc)
A Go client library for accessing a remote Selenium session and perform inspections. Complies W3C [specifications](https://www.w3.org/TR/webdriver/).

## Installation
Use the following command to download this module:
```
go get github.com/codedius/wdc
```

## Usage
```go
import "github.com/codedius/wdc"
```

Construct a new client:

```go
sess := wdc.Session{
    ID:  "bb6d-6bb9bd380a11",   // Webdriver session ID
    URL: "http://example.com",  // Webdriver remote server URL
}

client, err := wdc.New(&sess)
if err != nil {
    // error handling
}
```

Perform inspections. For example:


```go
ctx := context.Background()

// Navigate to the page
err = client.NavigateTo(ctx, "http://example.com") 
if err != nil {
    // error handling
}

// Set timeout to find the element
err = client.TimeoutElementFind(ctx, 10*time.Second)
if err != nil {
    // error handling
}

// Get element ID to perform further inspections
eid, err := client.ElementFind(ctx, wdc.ByXPath, "//div[@class=my-class]")
if err != nil {
    // error handling
}
	
// Get element's text
text, err := client.ElementText(ctx, eid)
if err != nil {
    // error handling
}
```

For more info please read library's [documentation](https://pkg.go.dev/github.com/codedius/wdc?tab=doc) from go.dev
