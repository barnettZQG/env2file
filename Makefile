
test-mysql: localbuild
	HOST_NAME=mysql11 MYSQLC_MYSQLD_PORT=3306 MYSQLC_MYSQLD_read_buffer_size=2M MYSQLC_MYSQLD_SERVER_ID=${HOST_NAME} MYSQL_MYSQLD_DATADIR=/data/cc build/env2file cre --path ./test/mysql/custom.cnf --format mysql	
test-redis: localbuild
	REDISC_PORT=6379 REDISC_TIMEOUT=0 build/env2file cre --path ./test/redis/custom.conf --format redis
localbuild:
	go build -o build/env2file
release:
	docker run --rm -it -v `pwd`:/go/src/github.com/barnettZQG/env2file -w /go/src/github.com/barnettZQG/env2file golang:1.11 go build -ldflags " -w" -o build/env2file-linux
	docker run --rm -e GOOS=windows -it -v `pwd`:/go/src/github.com/barnettZQG/env2file -w /go/src/github.com/barnettZQG/env2file golang:1.11 go build -ldflags " -w" -o build/env2file-win.exe
image: 
	docker build -t barnett/env2file .	
	