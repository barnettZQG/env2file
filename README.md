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