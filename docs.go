/*
go-rip is simple rest json api client. It has fluent chainable api, so you can do request along with unmarshalling
and error checking on single line.

First we prepare client that we will use in later examples:

	client := rip.New("http://localhost/api/v1").Header("Token", "Token").AppendSlash(true)

Now we can have a look at some examples of rest json api calls


	type Product struct {}
	product := Product{}

	status := 0

	if err := client.Get("product", 1).Do(context.Background(), &product).Status(&status); err != nil {
		panic(err)
	}

In this example we did following:

	- make GET http request to http://localhost/api/v1/product/1/
	- we have unmarshalled json response to `product` variable
	- we have set http status code into `status` variable
	- we have checked for error

This is awesome!

Majority of methods return Client so you can chain your calls. Structures are immutable, so chaining creates new struct.
Do method returns response, which provides also useful methods.

Let's do POST http call

	if err := client.Post("product").Data(product).Do(context.Background()).Status(&status).Error(); err != nil {
		panic(err)
	}

*/
package rip
