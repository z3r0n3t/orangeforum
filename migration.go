// Copyright (c) 2020 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"database/sql"
	"log"
)

func migrate0(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE domains (
		id INT NOT NULL AUTO_INCREMENT,
		domain_name VARCHAR(120) NOT NULL,
		forum_name VARCHAR(100) NOT NULL,
		signup_disabled BOOL NOT NULL DEFAULT 0,
		login_to_view BOOL NOT NULL DEFAULT 0,
		read_only BOOL NOT NULL DEFAULT 0,
		created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY(id),
		UNIQUE(domain_name),
		INDEX(domain_name)
	);`)
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}

	_, err = db.Exec(`CREATE TABLE config (
		id INT NOT NULL AUTO_INCREMENT,
		k VARCHAR(32) NOT NULL,
		v VARCHAR(120) NOT NULL,
		PRIMARY KEY(id),
		INDEX(k)
	);`)
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}

	_, err = db.Exec(`CREATE TABLE users (
		id INT NOT NULL AUTO_INCREMENT,
		domain_id INT NOT NULL,
		username VARCHAR(32) NOT NULL,
		passwdhash VARCHAR(250) NOT NULL,
		email VARCHAR(250) DEFAULT NULL,
		resetpass_key VARCHAR(50) DEFAULT NULL,
		is_admin BOOL NOT NULL DEFAULT 0,
		is_mod BOOL NOT NULL DEFAULT 0,
		is_banned BOOL NOT NULL DEFAULT 0,
		resetpass_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY(id),
		FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE,
		UNIQUE(domain_id, username),
		INDEX(domain_id, username),
		INDEX(domain_id, email),
		INDEX(domain_id, is_admin),
		INDEX(domain_id, is_mod),
		INDEX(domain_id, resetpass_key)
	);`)
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}

	_, err = db.Exec(`CREATE TABLE sessions (
		id INT NOT NULL AUTO_INCREMENT,
		domain_id INT NOT NULL,
		user_id INT NOT NULL,
		sess VARCHAR(64) NOT NULL,
		created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(id),
		FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		INDEX(domain_id, user_id),
		INDEX(domain_id, sess)
	);`)
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}

	_, err = db.Exec(`CREATE TABLE topics (
		id INT NOT NULL AUTO_INCREMENT,
		domain_id INT NOT NULL,
		user_id INT NOT NULL,
		title VARCHAR(250) NOT NULL,
		is_sticky BOOL NOT NULL DEFAULT 0,
		created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		n_comments INT NOT NULL,
		PRIMARY KEY(id),
		FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		INDEX(domain_id, is_sticky, created_date),
		INDEX(domain_id, is_sticky, updated_date),
		INDEX(domain_id, user_id, updated_date)
	);`)
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}

	_, err = db.Exec(`CREATE TABLE comments (
		id INT NOT NULL AUTO_INCREMENT,
		domain_id INT NOT NULL,
		topic_id INT NOT NULL,
		user_id INT NOT NULL,
		content TEXT NOT NULL,
		is_sticky BOOL NOT NULL DEFAULT 0,
		created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY(id),
		FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE,
		FOREIGN KEY(topic_id) REFERENCES topics(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		INDEX(domain_id, topic_id, is_sticky, created_date),
		INDEX(domain_id, user_id, updated_date)
	);`)
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}

	_, err = db.Exec(`INSERT INTO config(k, v) VALUES("version", "0");`)
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}
}

func migrate(db *sql.DB) {
	var ver string
	ver = "-1"
	db.QueryRow("SELECT v FROM config WHERE k=?;", "version").Scan(&ver)

	if ver == "-1" {
		log.Println("[INFO] Migrating DB to version 0.")
		migrate0(db)
	} else if ver == "0" {
		log.Fatalf("[ERROR] Database schema already up-to-date. No migration done.\n")
	}
}
