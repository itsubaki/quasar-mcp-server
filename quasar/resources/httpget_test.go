package resources_test

import (
	"fmt"

	"github.com/itsubaki/quasar-mcp-server/quasar/resources"
)

func ExampleHttpGet() {
	bytes, err := resources.HttpGet("https://github.com/itsubaki/quasar-mcp-server")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes)[:23])

	// Output:
	// <!DOCTYPE html>
}
