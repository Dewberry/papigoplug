package main

import (
	"os"

	plug "github.com/Dewberry/papigoplug/papigoplug"
)

func main() {
	plug.InitLog("info")

	allowedParams := plug.PluginParams{
		Required: []string{"first", "last"},
		Optional: []string{"middle"},
	}

	params, err := plug.ParseInput(os.Args, allowedParams)
	if err != nil {
		plug.Log.Fatal(err)
	}
	plug.Log.Infof("Params provided: %s", params)

	results := map[string]interface{}{"foo": "bar", "success": true}
	plug.PrintResults(results)
}
