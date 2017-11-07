# go-rip
[![Build Status](https://travis-ci.org/phonkee/go-rip.svg?branch=master)](https://travis-ci.org/phonkee/go-rip)
[![Coverage Status](https://coveralls.io/repos/github/phonkee/go-rip/badge.svg?branch=master)](https://coveralls.io/github/phonkee/go-rip?branch=master)

Small library for making json rest calls.
go-rip provides fluent chainable api that makes doing rest calls reslly simple.

## name
Why the name? Well `rest in peace` fits quite good as rest api client.

### Install

You can install this library simply by typing:

    go get github.com/phonkee/go-rip

### Examples:

Let's see this library in action.
First let's create a client that we will use in all examples

```go
client := rip.New("http://127.0.0.1/api/v1")
```

Client has even more methods to add:

* Headers - add default headers for all method calls
* AppendSlash - add trailing slash if not presented
* Base - set base url
* Client - set custom http.Client instead of default ones
* Data - add data to be sent as body, data can be:
    * string - is converted to []byte
    * byte[] - sent as is
    * other - will be json Marshalled to []byte
* Path - add path prefix
* QueryValues - add query values (query parameters)
* UserAgent - set custom user agent instead of default one

Now let's make a http GET request.

```go
type Product struct {}

product = Product{}
status := 0
if err := client.Get("product", 1).Do(context.Background(), &product).Status(&status).Error(); err != nil {
    panic("oops")
}
```

This makes a http GET call to "http://127.0.0.1/api/v1/product/1/" and unmarshals result
to `product` variable. We are also handling errors and filling status variable with http status code in the same line.
Isn't that nice?

You can see that go-rip uses golangs context for http requests and must be passed as first argument to `Do` method.

Let's make a POST request with some data.

```go
product = Product{}
result = map[string]interface{}
if err := client.Post("product", 1).Data(product).Do(context.Background(), &result).Error(); err != nil {
    panic(err)
}
```

We just made http POST request to "http://127.0.0.1/api/v1/product/1/" and we sent body of the request
json marshalled product. Ain't that really nice?

### author
Peter Vrba <phonkee@phonkee.eu>
