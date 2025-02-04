# http-server

## Installation

```bash
go install github.com/snowdreamtech/go-http-server
```

## Usage

```bash
$ ./http-server --help
A Simple Static HTTP Server built with gin and golang.

Usage:
   [flags]
   [command]       

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  env         Print  version and environment info
  help        Help about any command
  version     Print the version number of

Flags:
      --autoindex-exact-size           For the HTML format, specifies whether exact file sizes should be output in the directory listing,
                                       or rather rounded to kilobytes, megabytes, and gigabytes.
      --autoindex-time-format string   this is the AutoIndex Time Format. (default "2006-01-02 15:04:05")
      --basic                          If it is set, we will use HTTP Basic authentication.

                                       Used together with -u, --user <user:password>.

                                       Providing --basic multiple times has no extra effect.

                                       Example:-u name:password --basic https://example.com
  -c, --config string                  If it is not set, we will try with development.(json/env/ini/yaml/toml/hcl/properties
      --contact-email string           HTTPS Contact Email
      --enable-https                   If it is set, we will enable https.
  -g, --gzip                           If it is set, we will compress with gzip. (default true)
  -h, --help                           help for this command
  -H, --host string                    Host optionally specifies the Http Address for the server to listen on,
                                       in the form "host:port". If empty, "host:" (host localhost) is used.
      --https-cert-dir string          HTTPS Cert Directory (default "certs")
      --https-cert-file string         HTTPS Cert File
      --https-domains stringArray      HTTPS Domains
      --https-key-file string          HTTPS Key File
      --https-port string              HTTPS PORT
      --log-dir string                 The Log Directory which store access.log, error.log etc. (default ".")
  -P, --port string                    Port optionally specifies the TCP Port for the server to listen on,
                                       in the form "host:port". If empty, ":port" (port 8080) is used.
      --preview-html                   For static web files, Whether preview them. (default true)
      --rate-limiter string            Define a limit rate to several requests per hour.
                                        You can also use the simplified format "<limit>-<period>"", with the given
                                        periods:

                                        * "S": second
                                        * "M": minute
                                        * "H": hour
                                        * "D": day

                                        Examples:

                                        * 5 reqs/second: "5-S"
                                        * 10 reqs/minute: "10-M"
                                        * 1000 reqs/hour: "1000-H"
                                        * 2000 reqs/day: "2000-D"
                                         (default "5-S")
      --read-timeout int               ReadTimeout is the maximum duration for reading the entire
                                       request, including the body. A zero or negative value means
                                       there will be no timeout.

                                       Because ReadTimeout does not let Handlers make per-request
                                       decisions on each request body's acceptable deadline or
                                       upload rate, most users will prefer to use
                                       ReadHeaderTimeout. It is valid to use them both. (default 10)
      --referer-limiter                Limit by referer
      --speed-limiter int               Specify  the  maximum  transfer  rate you want curl to use - for
                                       downloads.
                                       The given speed is measured in bytes/second,
  -u, --user string                    Specify the user name and password to use for server authentication.

                                       The user name and passwords are split up  on  the  first  colon,
                                       which  makes  it impossible to use a colon in the user name with
                                       this option. The password can, still. (default "admin:admin")
      --write-timeout int              WriteTimeout is the maximum duration before timing out
                                       writes of the response. It is reset whenever a new
                                       request's header is read. Like ReadTimeout, it does not
                                       let Handlers make decisions on a per-request basis.
                                       A zero or negative value means there will be no timeout. (default 10)
  -w, --wwwroot string                 By default, the wwwroot folder is treated as a web root folder.
                                       Static files can be stored in any folder under the web root and accessed with a relative path to that root.

Use " [command] --help" for more information about a command.

```

## License 

```bash
(The MIT License)

Copyright (c) 2023-present SnowdreamTech Inc.

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```