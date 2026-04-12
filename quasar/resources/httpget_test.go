package resources_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/itsubaki/quasar-mcp-server/quasar/resources"
)

func ExampleHttpGet() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "test response body")
	}))
	defer srv.Close()

	bytes, err := resources.HttpGet(srv.URL)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	// Output:
	// test response body
}
