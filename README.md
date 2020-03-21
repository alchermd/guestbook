# Guestbook 

An implementation of a guestbook app with Go and MySQL

## Setup

This project depends on the following dependencies:

```bash
$ go get github.com/gorilla/mux
$ go get github.com/go-sql-driver/mysql
```

Then run the initial database schema:

```bash
$ mysql < init.sql
```

Lastly, set the environment variables for database connection:

```bash
$ export DB_USER=root
$ export DB_PASSWORD=yourpasswordhere
$ go run app.go
2020/03/21 23:44:46 Serving static assets at /static
2020/03/21 23:44:46 Starting server at port 3000
```

## License

Released under [MIT](LICENSE)