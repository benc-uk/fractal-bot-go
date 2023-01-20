package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	customHandlerPort := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if customHandlerPort == "" {
		customHandlerPort = "8080"
	}

	mux := http.NewServeMux()

	// This handler accepts standard HTTP requests from the Functions runtime and returns a vanilla HTTP response
	// Due to having a HTTP binding and enableForwardingHttpRequest set to true
	mux.HandleFunc("/api/fractal", fractalHandler)

	// This handler is has a timer binding, so is invoked by the Functions runtime with a special payload
	// And is expected to return JSON in a specific format
	// See https://learn.microsoft.com/en-us/azure/azure-functions/functions-custom-handlers#request-payload
	mux.HandleFunc("/tweetFractal", tweetHandler)

	log.Printf("### 🚀 Starting custom Go handler")
	log.Printf("### 🌐 Server listening on: %s", customHandlerPort)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
}
