# Assessment

If your GOPATH is setup correctly you should just have to run `go get github.com/jmlattanzi/golang-assessment`, navigate into the directory and run `go run main.go` or if you'd like to build it and execute it, run `go build`

This project uses an external PostgreSQL database. Connection to this DB is made through a file called `config.json` in the root of this project. If you would like to use your own PostgreSQL database, put the connection url in the `config.json` with the key `DBURL`, or contact me for the config file. You can email me at jonathanlattanzi@gmail.com
