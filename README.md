## 简介
 由环境变量生成配置文件的命令行工具，常用于集成的Docker容器中，容器启动后、服务真正进程启动前执行。

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
custom.conf:
```
port 6379
```
