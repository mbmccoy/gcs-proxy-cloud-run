// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/DomZippilli/gcs-proxy-cloud-function/main/common"
	"github.com/DomZippilli/gcs-proxy-cloud-function/main/proxyhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	storage "cloud.google.com/go/storage"
)

func setup() {
	// configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// set the bucket name from environment variable
	common.BUCKET = os.Getenv("BUCKET_NAME")

	// initialize the client
	c, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal().Msgf("main: %v", err)
	}
	common.GCS = c
}

func main() {
	log.Print("starting server...")
	setup()
	http.HandleFunc("/", ProxyHTTPGCS)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Warn().Msgf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal().Msgf("main: %v", err)
	}
}

// ProxyHTTPGCS is the entry point for the cloud function, providing a proxy that
// permits HTTP protocol usage of a GCS bucket's contents.
func ProxyHTTPGCS(output http.ResponseWriter, input *http.Request) {
	ctx := context.Background()
	// route HTTP methods to appropriate handlers.
	// ===> Your filters go below here <===
	switch input.Method {
	case http.MethodGet:
		proxyhttp.GET(ctx, output, input)
	default:
		http.Error(output, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
	return
}
