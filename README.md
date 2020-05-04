Orange Forum
============

Note: Orange Forum 2.0 is work in progress. Please see [orangeforum-1.x.x branch](https://github.com/s-gv/orangeforum/tree/orangeforum-1.x.x) for the latest stable version.

[Orange Forum](http://www.goodoldweb.com/orangeforum/) is an easy to deploy forum that has minimal dependencies and uses very little javascript. It is written is golang and a [compiled binary](https://github.com/s-gv/orangeforum/releases) is available for linux. Try the latest version hosted [here](https://groups.goodoldweb.com/). Please contact [info@goodoldweb.com](mailto:info@goodoldweb.com) if you have any questions or want support.

How to use
----------

By default, sqlite is used, so it's easy to get started.
[Download](https://github.com/s-gv/orangeforum/releases) the binary and migrate the database with:

```
./orangeforum -migrate
```

Create an admin:

```
./orangeforum -createadmin
```

Finally, start the server:

```
./orangeforum
```

Notes
-----

There are three types of users in Orangeforum: admin, mod, and regular users. Admins are the most previleged and can do anything. Mods can edit posts and ban users.

Dependencies
------------

- Go 1.8 (only for compiling)
- Postgres 9.5 (or use embedded sqlite3)

Options
-------

- `-addr <port>`: Use `./orangeforum -addr 8086` to listen on port 8086.
- `-dbdriver <db>` and `-dsn <data_source_name>`: PostgreSQL and SQLite are supported. SQLite is the default driver.

To use postgres, run `./orangeforum -dbdriver postgres -dsn postgres://pguser:pgpasswd@localhost/orangeforum`

To save an sqlite db at a different location, run `./orangeforum -dsn path/to/myforum.db`.

Commands
--------

- `-help`: Show a list of all commands and options.
- `-migrate`: Migrate the database. Run this once after updating the orangeforum binary (or when starting afresh).
- `-createadmin`: Create an admin.
- `-createuser`: Create a new user with no special privileges.
- `-changepasswd`: Change password of a user.
- `-deletesessions`: Drop all sessions and log out all users.

