// Copyright (c) 2019 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr := flag.String("addr", ":9123", "Port to listen on")

	mux := http.NewServeMux()

	srv := &http.Server{
		Handler:      mux,
		Addr:         *addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("[INFO] Starting orangeforum at", *addr)

	err := srv.ListenAndServe()
	if err != nil {
		log.Panicf("[ERROR] %s\n", err)
	}
}
