# Go WebDriver Remote Client
A Go client library for accessing a remote webdriver session and perform page inspections. Partially (see the full list of implemented members in the [doc](https://pkg.go.dev/github.com/codedius/wdc?tab=doc)) complies W3C [specifications](https://www.w3.org/TR/webdriver/). Library is intended to be as simple as possible and not claiming to be full-fledged webdriver client. If you’re looking for that kind of client you may try another awsome project [selenium-go](https://github.com/tebeka/selenium).

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
    ID:  "ba97d362c7880e1adc7b81e613c665c6",   // webdriver session ID
    URL: "http://example.com",  	       // webdriver remote server URL
}

client, err := wdc.New(&sess)
if err != nil {
    // error handling
}
```

Perform inspections. For example:


```go
ctx := context.Background()

// navigate to the page
err = client.NavigateTo(ctx, "http://example.com") 
if err != nil {
    // error handling
}

// set timeout to find the element
err = client.TimeoutElementFind(ctx, 10*time.Second)
if err != nil {
    // error handling
}

// get element's ID to perform further inspections
eid, err := client.ElementFind(ctx, wdc.ByXPath, "//div[@class='my-class']")
if err != nil {
    // error handling
}
	
// get element's text
text, err := client.ElementText(ctx, eid)
if err != nil {
    // error handling
}
```
