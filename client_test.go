package rip

import (
	"testing"

	"net/http"

	"context"

	"net/url"

	. "github.com/smartystreets/goconvey/convey"
)

// internalClient returns internal implementation of Client
func internalClient(c Client) client {
	return c.(client)
}

func TestClient(t *testing.T) {

	Convey("Test URL", t, func() {

		prefix := "/api/v1/"
		baseUrl := "http://127.0.0.1" + prefix

		data := []struct {
			path        []interface{}
			resultPath  string
			appendSlash bool
		}{
			{
				[]interface{}{"this", "is", "product", 1},
				prefix + "this/is/product/1",
				false,
			},
			{
				[]interface{}{"this", "is", "product", 1},
				prefix + "this/is/product/1/",
				true,
			},
			{
				[]interface{}{"this", "is", "", "", " ", "product", 1},
				prefix + "this/is/product/1/",
				true,
			},
		}

		for _, item := range data {
			client := New(baseUrl).Path(item.path...).AppendSlash(item.appendSlash)
			So(client.URL().Path, ShouldEqual, item.resultPath)
		}
	})

	Convey("Test Query Values", t, func() {
		client := New("http://127.0.0.1/api/v1").AppendSlash(true).QueryValues(url.Values{
			"key":  []string{"value"},
			"oops": []string{"oops"},
		})
		So(client.URL().String(), ShouldEqual, "http://127.0.0.1/api/v1/?key=value&oops=oops")

		client2 := New("http://127.0.0.1/api/v1").AppendSlash(true).QueryValues(url.Values{
			"key":  []string{"value"},
			"oops": []string{"oops"},
		}).QueryValues(nil)

		So(client2.URL().String(), ShouldEqual, "http://127.0.0.1/api/v1/")
	})

	Convey("Test Method", t, func() {
		data := []string{
			http.MethodDelete,
			http.MethodGet,
			http.MethodHead,
			http.MethodPatch,
			http.MethodPost,
			http.MethodPut,
		}

		for _, method := range data {
			So(New().Method(method).Request().Method, ShouldEqual, method)
		}
	})

	Convey("Test UserAgent", t, func() {
		data := []string{
			"user-agent",
			"test-agent",
		}

		for _, ua := range data {
			So(New().UserAgent(ua).Request().Header.Get("User-Agent"), ShouldEqual, ua)
		}

	})

	Convey("Test HTTP methods", t, func() {

		client := New()
		data := []struct {
			method     func(parts ...interface{}) Client
			httpMethod string
		}{
			{client.Delete, http.MethodDelete},
			{client.Get, http.MethodGet},
			{client.Head, http.MethodHead},
			{client.Options, http.MethodOptions},
			{client.Patch, http.MethodPatch},
			{client.Post, http.MethodPost},
			{client.Put, http.MethodPut},
		}

		for _, item := range data {
			So(item.method().Request().Method, ShouldEqual, item.httpMethod)
		}
	})

	Convey("Test Get", t, func() {
		client := New().Get()
		So(internalClient(client).method, ShouldEqual, http.MethodGet)
	})

	Convey("Test Request", t, func() {
		client := New().Delete().Header("key", "value")
		request := client.Request()
		So(request, ShouldNotBeNil)

		So(request.Method, ShouldEqual, http.MethodDelete)
		So(request.URL, ShouldResemble, client.URL())
		So(request.Header.Get("key"), ShouldEqual, "value")
	})

	Convey("Test Do valid", t, func() {
		client := New("http://dog.ceo/api/")

		target := map[string]interface{}{}

		var status int

		err := client.Get("breeds/list/all").Do(context.Background(), &target).Status(&status).Error()
		So(err, ShouldBeNil)
		So(status, ShouldEqual, http.StatusOK)
		So(target["status"].(string), ShouldEqual, "success")

	})

	Convey("Test Do invalid", t, func() {
		client := New("asdf")

		target := map[string]interface{}{}

		var status int

		err := client.Get("breeds/list/all").Do(context.Background(), &target).Status(&status).Error()
		So(err, ShouldNotBeNil)
		So(status, ShouldEqual, 0)
	})

}
