[![Build Status](https://travis-ci.org/barnettZQG/env2file.svg?branch=master)](https://travis-ci.org/barnettZQG/env2file)
## 简介
 由环境变量生成配置文件的命令行工具，常用于集成的Docker容器中，容器启动后、服务真正进程启动前执行。

## 安装

```bash
wget https://github.com/barnettZQG/env2file/releases/download/0.1.1/env2file-linux
```

docker环境
```
docker pull barnett/env2file
# test
docker run -it --rm barnett/env2file env2file
```
## 使用说明

```
NAME:
   env2file - A new cli application

USAGE:
   env2file [global options] command [command options] [arguments...]

COMMANDS:
     conversion, conv  Renders the specified file template by reading the environment variable
     create, cre       Generates a configuration file to the specified path by reading the environment variable and following the specified format
     help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## 生成指定格式的配置文件
`env2file create` 命令根据指定规范生成配置文件,现支持 mysql redis and default
### 关于生成mysql配置文件
读取以`MYSQLC_`开头的环境变量，以下划线分割，第二部分为配置块名称，比如mysqld 或者client. 第三部分为配置名
例如： MYSQLC_MYSQLD_PORT = 3306
custom.cnf:
```
[mysqld]
 port = 3306
```
### 关于生成redis配置文件
读取以`REDISC_`开头的环境变量，以下划线分割,第二部分为配置名
例如： REDISC_PORT = 6379
命令：
```
env2file create --path redis.conf --perm 0755 --format redis
```

custom.conf:
```
port 6379
```

### 关于生成默认配置文件
读取以`C_`开头的环境变量，以下划线分割,第二部分为配置名
例如： C_MYNAME = barnett
custom.conf:
```
myname=barnett
```
以上类型均支持值热渲染,即如果在环境变量值中具有 ${XXX} 的规范字符串，将继续渲染值，使用对应的环境变量，仅支持一级渲染。

## 模版配置文件重写
`env2file conversion` 命令根据指定的配置文件模版，需要渲染的变量使用 ${XXX} 的形式定义。工具将根据从环境变量中获取值来渲染并重写文件，仅支持一级渲染。
 * 定义默认值
 ${XXX:123} 此方式定义的变量若从环境变量中无法读取值时，将使用`123`作为默认值
