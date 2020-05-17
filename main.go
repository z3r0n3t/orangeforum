// Copyright (c) 2020 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var er error
	var db *sql.DB

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr := flag.String("addr", ":9123", "Port to listen on")
	dsn := flag.String("dsn", "mysqluser:mysqlpasswd@tcp(10.156.14.46)/testmysql", "Database source name (ex: mysqluser:mysqlpasswd@tcp(127.0.0.1)/orangeforum)")
	doMigrate := flag.Bool("migrate", true, "Migrate database to current version")

	db, er = sql.Open("mysql", *dsn)
	if er != nil {
		log.Panicf("[ERROR] Error opening database: %s\n", er)
	}

	if *doMigrate {
		migrate(db)
		os.Exit(0)
	}

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
