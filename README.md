# jedi

`jedi` is a basic multi-resource searching matcher. 

The purpose of `jedi` is searching matched data from different types of the given data file. `jedi` invokes network syscall and encoding package to decode xml/json data into structured data, and `jedi` also use goroutine to ensure the speed of these actions.

## Feature

Please visit todo list.

## Roadmap

First stage, jedi will implement basic function, make sure this stuff avaiable. Second stage, we will improve jedi to showing error message and log clearly, probably support cmd runing and cron task. Last stage, we'll expand jedi to support more middlewares, just like Redis/RPC etc., we want to make it be a stronger network searching matcher. I mean, this is our perspective, we'll do it well, or not.

## Getting Started

### Usage

```shell
» ./jedi -h
Usage of ./jedi:
  -keyword string
    	specific a keyword to query (default "China")
  -limit int
    	how many matched items shows here (default 20)

» ./jedi -keyword="social" -limit=10
```


## TODO

- [x] support concurrency search
- [x] godoc support
- [x] enhenced logger
- [x] test file fully coveraged
- [x] add process bar support
- [x] docker integraty for logger
- [x] support output message color
- [x] fix code format, repalce `if err != nil` with noerr
- [x] cmd and argument control support
- [ ] support different type data file(like xml etc.)
- [ ] imporve performance(replace json with jsoniter)
- [ ] generic support
- [ ] data store: MySQL/PostgreSQL 
- [ ] cache results in local database: Redis/PostgreSQL
- [ ] service discovery: etcd/etcd-proxy
- [ ] cron task 


## Architecture

![jedi.jpg](https://github.com/i0Ek3/jedi/blob/master/drawio/jedi.jpg)


## Coding Rule

Please check [this post](https://github.com/Tencent/secguide/blob/main/Go安全指南.md).


## Q&A

> 1. When I use command `make gdb` to debug the program, it show me error message like this: `Unable to find Mach task port for process-id 22182: (os/kern) failure (0x5). (please check gdb is codesigned - see taskgated(8))`, how can I fix that?

Check [this post](https://gist.github.com/gravitylow/fb595186ce6068537a6e9da6d8b5b96d).

ps: Enable crsutil will conflict with proxychains4-ng, so make your choice.


## Contributing

PRs and Issues are also welcome.


## Credit

- goinaction


## License

MIT.
