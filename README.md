# go-rip

Small library for making json rest calls. Why the name? Well `rest in peace` fits quite good as name.

### Install

You can install this library simply by typing:

    go get github.com/phonkee/go-rip


### Examples:

Let's see this library in action.
First let's create a client that we will use in all examples

```go
client := rip.New("http://127.0.0.1/api/v1").AppendSlash(true).Header("Token", "Token")
```

You can see that api is really fluent and chainable.
Now let's make a http GET request.

```go
type Product struct {}

product = Product{}

if err := client.Get("product", 1).Do(context.Background(), &product); err != nil {
    panic("oops")
}
```

This makes a rest api call to "http://127.0.0.1/api/v1/product/1/" and unmarshals result
to `product` variable. We are also handling errors in the same line. Isn't that nice?

Let's make a post request with some data.

```go
product = Product{}
result = map[string]interface{}
if err := client.Post("product", 1).Data(product).Do(context.Background(), &result); err != nil {
    panic(err)
}
```

We just made http POST request to "http://127.0.0.1/api/v1/product/1/" and we sent body of the request
json marshalled product. Ain't that really nice?

### author
Peter Vrba <phonkee@phonkee.eu>