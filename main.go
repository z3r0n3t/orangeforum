// Copyright (c) 2020 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"bufio"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println("")

	fmt.Print("Enter Password (again): ")
	bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
	password2 := string(bytePassword2)
	fmt.Println("")

	if password != password2 {
		log.Fatalf("[ERROR] Passwords do not match")
	}

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func main() {
	var er error
	var db *sql.DB

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr := flag.String("port", ":9123", "Port to listen on")
	dsn := flag.String("dsn", "mysqluser:mysqlpasswd@tcp(10.156.14.46)/testmysql", "Database source name")
	doMigrate := flag.Bool("migrate", false, "Migrate database schema to current version")
	createDomain := flag.Bool("create-domain", false, "Create new domain")
	renameDomain := flag.Bool("rename-domain", false, "Rename existing domain")
	createUser := flag.Bool("create-user", false, "Create new user")
	createAdminUser := flag.Bool("create-admin-user", false, "Create new admin user")
	changePasswd := flag.Bool("change-password", false, "Change password")
	banUser := flag.Bool("ban-user", false, "Ban user")
	unbanUser := flag.Bool("unban-user", false, "Unban user")
	dropSessions := flag.Bool("drop-sessions", false, "Drop sessions and log out all users")
	enableReadOnly := flag.Bool("enable-readonly", false, "Enable read-only mode")
	disableReadOnly := flag.Bool("disable-readonly", false, "Disable read-only mode")

	flag.Parse()

	db, er = sql.Open("mysql", *dsn)
	if er != nil {
		log.Fatalf("[ERROR] Error opening database: %s\n", er)
	}

	if *doMigrate {
		migrate(db)
		return
	}

	var ver string
	er = db.QueryRow("SELECT v FROM config WHERE k='version';").Scan(&ver)
	if er != nil {
		log.Fatalf("[ERROR] Database migration may be needed. Error reading schema version: %s\n", er)
	}
	if ver != "0" {
		log.Fatalf("[ERROR] Database migration may be needed. Incorrect schema version -- Expected: 0, Got: %s\n", ver)
	}

	if *createDomain {
		fmt.Print("Enter domain name (ex: www.google.com): ")
		var domainName string
		fmt.Scanln(&domainName)

		forumName := "Orange Forum"

		_, er = db.Exec("INSERT INTO domains(domain_name, forum_name) VALUES(?, ?);", domainName, forumName)
		if er != nil {
			log.Fatalf("[ERROR] Error creating domain: %s\n", er)
		}
		return
	}

	if *renameDomain {
		fmt.Print("Enter current domain name (ex: www.google.com): ")
		var domainName string
		fmt.Scanln(&domainName)

		var domainID int
		er = db.QueryRow("SELECT id FROM domains WHERE domain_name=?;", domainName).Scan(&domainID)
		if er != nil {
			log.Fatalf("[ERROR] Could not find domain: %s\n", domainName)
		}

		fmt.Print("Enter new domain name: ")
		var newDomainName string
		fmt.Scanln(&newDomainName)

		_, er = db.Exec("UPDATE domains SET domain_name=? WHERE id=?;", newDomainName, domainID)
		if er != nil {
			log.Fatalf("[ERROR] Error renaming domain: %s\n", er)
		}
		return
	}

	if *createUser || *createAdminUser {
		fmt.Print("Enter domain name (ex: www.google.com): ")
		var domainName string
		fmt.Scanln(&domainName)

		var domainID int
		er = db.QueryRow("SELECT id FROM domains WHERE domain_name=?;", domainName).Scan(&domainID)
		if er != nil {
			log.Fatalf("[ERROR] Could not find domain: %s\n", domainName)
		}

		username, passwd := credentials()

		if passwd == "" {
			log.Fatalf("[ERROR] Password cannot be blank\n")
		}

		if passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err == nil {
			if *createAdminUser {
				_, er = db.Exec("INSERT INTO users(domain_id, username, passwdhash, is_admin) VALUES(?, ?, ?, 1);", domainID, username, hex.EncodeToString(passwdHash))
			} else {
				_, er = db.Exec("INSERT INTO users(domain_id, username, passwdhash) VALUES(?, ?, ?);", domainID, username, hex.EncodeToString(passwdHash))
			}
		} else {
			log.Fatalf("[ERROR] Error hashing password: %s\n", err)
		}
		return
	}

	if *changePasswd {
		fmt.Print("Enter domain name (ex: www.google.com): ")
		var domainName string
		fmt.Scanln(&domainName)

		var domainID int
		er = db.QueryRow("SELECT id FROM domains WHERE domain_name=?;", domainName).Scan(&domainID)
		if er != nil {
			log.Fatalf("[ERROR] Could not find domain: %s\n", domainName)
		}

		userName, passwd := credentials()

		var userID int
		er = db.QueryRow("SELECT id FROM users WHERE domain_id=? AND username=?;", domainID, userName).Scan(&userID)
		if er != nil {
			log.Fatalf("[ERROR] Could not find user: %s\n", userName)
		}

		if passwd == "" {
			log.Fatalf("[ERROR] Password cannot be blank\n")
		}

		if passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err == nil {
			db.Exec("UPDATE users SET passwdhash=? WHERE domain_id=? AND id=?;", hex.EncodeToString(passwdHash), domainID, userID)
		} else {
			log.Fatalf("[ERROR] Error hashing password: %s\n", err)
		}
		return
	}

	if *banUser || *unbanUser {
		fmt.Print("Enter domain name (ex: www.google.com): ")
		var domainName string
		fmt.Scanln(&domainName)

		var domainID int
		er = db.QueryRow("SELECT id FROM domains WHERE domain_name=?;", domainName).Scan(&domainID)
		if er != nil {
			log.Fatalf("[ERROR] Could not find domain: %s\n", domainName)
		}

		fmt.Print("Enter username: ")
		var userName string
		fmt.Scanln(&userName)

		var userID int
		er = db.QueryRow("SELECT id FROM users WHERE domain_id=? AND username=?;", domainID, userName).Scan(&userID)
		if er != nil {
			log.Fatalf("[ERROR] Could not find user: %s\n", userName)
		}

		if *banUser {
			db.Exec("UPDATE users SET is_banned=? WHERE id=?;", true, userID)
			db.Exec("DELETE FROM sessions WHERE domain_id=? AND user_id=?;", domainID, userID)
		} else {
			db.Exec("UPDATE users SET is_banned=? WHERE id=?;", false, userID)
		}
		return
	}

	if *enableReadOnly || *disableReadOnly {
		fmt.Print("Enter domain name (ex: www.google.com): ")
		var domainName string
		fmt.Scanln(&domainName)

		var domainID int
		er = db.QueryRow("SELECT id FROM domains WHERE domain_name=?;", domainName).Scan(&domainID)
		if er != nil {
			log.Fatalf("[ERROR] Could not find domain: %s\n", domainName)
		}

		if *enableReadOnly {
			db.Exec("UPDATE domains SET read_only=? WHERE id=?;", true, domainID)
		} else {
			db.Exec("UPDATE domains SET read_only=? WHERE id=?;", false, domainID)
		}
		return
	}

	if *dropSessions {
		fmt.Print("Enter domain name (leave blank to drop sessions for all domains): ")
		var domainName string
		fmt.Scanln(&domainName)

		if domainName == "" {
			db.Exec("DELETE FROM sessions;")
		} else {
			var domainID int
			er = db.QueryRow("SELECT id FROM domains WHERE domain_name=?;", domainName).Scan(&domainID)
			if er != nil {
				log.Fatalf("[ERROR] Could not find domain: %s\n", domainName)
			}
			db.Exec("DELETE FROM sessions WHERE id=?;", domainID)
		}
		return
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
		log.Fatalf("[ERROR] %s\n", err)
	}
}
