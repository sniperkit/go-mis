MIS 2.0 - Micro-Services using Go

# Deploy Branch
Deploy branch is connected to jenkins. It will build and deploy when you push to deploy branch

# PRE-REQUISITE

Create a `config.json` file with below JSON sample.

```
{
	"psql": [
		{
			"db": "mis_amartha_dev",
			"host": "localhost",
			"port": "32003",
			"username": "postgres",
	    "password": "postgres",
			"ssl_mode": "disable"
		},
		{
			"db": "survey",
			"host": "localhost",
			"port": "32003",
			"username": "postgres",
	    "password": "postgres",
			"ssl_mode": "disable"
		}
	],

  "mysql": {
		"db": "db_mis_amartha",
		"host": "localhost",
		"port": "32003",
		"username": "root",
    "password": null,
		"ssl_mode": "disable"
	}
}
```

# HOW TO RUN

```
$ go run main.go -h
```

# DEPENDENCIES

```
$ go get -u gopkg.in/kataras/iris.v4
$ go get -u gopkg.in/iris-contrib/middleware.v4/secure
$ go get -u gopkg.in/iris-contrib/middleware.v4/recovery
$ go get -u gopkg.in/iris-contrib/middleware.v4/logger
$ go get -u github.com/jinzhu/gorm
$ go get -u github.com/jinzhu/gorm/dialects/postgres
$ go get -u github.com/nferruzzi/gormGIS
$ go get -u github.com/parnurzeal/gorequest
```
