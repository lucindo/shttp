# shttp

[![Build Status](https://drone.io/github.com/lucindo/shttp/status.png)](https://drone.io/github.com/lucindo/shttp/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucindo/shttp)](https://goreportcard.com/report/github.com/lucindo/shttp)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.png?v=103)](https://opensource.org/licenses/mit-license.php)

`shttp` is a zero configuration HTTP server to help on development of web apps. After having problems with some alternatives like `python -m SimpleHTTPServer` I decided to write this tool.

There are 3 main operation modes:

 1. Directory: expose a local directory via HTTP (default mode).
 2. Proxy: act as a reverse proxy to another server.
 3. API: return `200 OK` to all requests with good logs for debugging.

### Get `shttp`

`shttp` is a Golang program and you can get the binary for your platform here: [Releases](https://github.com/lucindo/shttp/releases).

If you have Go installed you may use:

```sh
go get -u github.com/lucindo/shttp
```

### Examples of use

To start a HTTP server on current directory (default bind address is `localhost:8080`):

```sh
shttp
```

Some options are enabled by default: CORS support and HTTP cache headers setted to 0 (no cache). It will log one line for request using the Apache Common Log format.

You can set an directory using the `-dir` option. One handy flag is `-open` that open a browser pointing to the server:

```sh
shttp -dir /path/to/my/webapp -open
```

The default file is `index.html`.

If you're writting some API client you can start `shttp` this way:

```sh
shttp -api -debug
```

The server you only give `200 OK` responses for all requests. The `debug` option you show you detailed information about your request (protocol, path, query string parameters, form parameters, all headers and so on).

`shttp` can also act as a reverse proxy to another server:

```sh
shttp -proxy http://my.cool.server/
```

This mode is very useful in conjunction of `-debug` flag.

### Options

Current options are:

```
Usage of shttp:
  -api
    	Catch-all handler to debug requests
  -cache
    	Don't add headers disabling HTTP cache for all requests
  -debug
    	Log request information
  -dir string
    	Directory to expose (default ".")
  -host string
    	Listen address (default "localhost")
  -maxheaders int
    	Max header size in bytes (default 1048576)
  -nocors
    	Disable CORS headers
  -open
    	Open a browser pointing to this server
  -port int
    	Port to bind the server (default 8080)
  -proxy string
    	Act as a reverse proxy
  -quiet
    	Do not log requests
  -rtimeout duration
    	Server read timeout (default 10s)
  -wtimeout duration
    	Server write timeout (default 10s)
```

### ToDo

- [ ] SSL support (maybe using Let's Encrypy)
- [ ] Log to file
- [ ] Daemon mode
