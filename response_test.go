package rip

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"context"
	"fmt"
)

func TestResponse(t *testing.T) {

	Convey("Test New", t, func() {
		Convey("Test Something", func() {
			client := New("http://dog.ceo/api/").Get("breeds/list/all")
			response := client.Do(context.Background())
			So(fmt.Sprintf("%T", client), ShouldEqual, fmt.Sprintf("%T", response.Client()))

			So(response.Header(), ShouldNotBeNil)
			So(response.Do(context.Background()).Body(), ShouldResemble, response.Body())

			var into []byte
			response.Raw(&into)
			So(len(into), ShouldBeGreaterThan, 0)

		})
	})
}
