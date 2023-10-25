package main

import (
	"fmt"
	"os"

	"github.com/Dewberry/papigoplug/papigoplug"
)

func main() {
	allowedParams := papigoplug.PluginParams{
		Required: []string{"first", "last"},
		Optional: []string{"middle"},
	}

	params, err := papigoplug.ParseInput(os.Args, allowedParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Params provided: %s\n", params)

	results := map[string]interface{}{"foo": "bar", "success": true}
	papigoplug.PrintResults(results)
}
