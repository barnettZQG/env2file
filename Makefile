build:
	go build -o build/env2file
test-mysql:build
	MYSQL_MYSQLD_PORT=3306 MYSQL_MYSQLD_DATADIR=/data/cc build/env2file cre --path ./test/mysql/custom.cnf --format mysql	
test-redis:build
	REDIS_PORT=6379 REDIS_TIMEOUT=0 build/env2file cre --path ./test/redis/custom.conf --format redis
release:
	docker run --rm -it -v `pwd`:/go/src/github.com/barnettZQG/env2file -w /go/src/github.com/barnettZQG/env2file golang:1.11 go build -ldflags " -w" -o build/env2file-linux
	docker run --rm -e GOOS=windows -it -v `pwd`:/go/src/github.com/barnettZQG/env2file -w /go/src/github.com/barnettZQG/env2file golang:1.11 go build -ldflags " -w" -o build/env2file-win.exe
	