# WebDriver Client (wdc)
A Go client library for accessing a remote Selenium session and perform inspections. Partially (see full list of implemented members below) Complies W3C [specifications](https://www.w3.org/TR/webdriver/). Library is intended to be as simple as possible and not claiming to be full-fledged Selenium client. If you’re looking for that kind of client you can try another awsome project [selenium-go](https://github.com/tebeka/selenium).

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

For more info please read library's [documentation](https://pkg.go.dev/github.com/codedius/wdc?tab=doc) from go.dev.

## Implemented members

