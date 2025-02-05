package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/crypto/acme/autocert"
	"snowdream.tech/gserver/middlewares"
	"snowdream.tech/gserver/pkg/configs"
	"snowdream.tech/gserver/pkg/env"
	glog "snowdream.tech/gserver/pkg/log"
	gnet "snowdream.tech/gserver/pkg/net"
	ghttp "snowdream.tech/gserver/pkg/net/http"
	ghttps "snowdream.tech/gserver/pkg/net/https"
	"snowdream.tech/gserver/pkg/tools"
)

var (
	rootCmd  *cobra.Command
	gHandler *gin.Engine
)

func init() {
	// Disable automaxprocs log
	// https://github.com/uber-go/automaxprocs/issues/19#issuecomment-557382150
	nopLog := func(string, ...interface{}) {}
	maxprocs.Set(maxprocs.Logger(nopLog))

	//load configs from File
	conf := configs.InitConfig()

	rootCmd = &cobra.Command{
		Use:   env.ProjectName,
		Short: env.ProjectName + " is a simple static http server",
		Long:  "A Simple Static HTTP Server built with gin and golang.",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			// Load configs from Cli
			if conf == nil && *configs.ConfigFile() != "" {
				conf = configs.InitConfig()
			}

			app := configs.GetAppConfig()

			if args != nil && len(args) == 1 {
				app.WwwRoot = args[0]
			}

			glog.InitLoggerConfig()

			welcome()

			tools.DebugPrintF("[INFO] Starting Web Server %s", env.ProjectName)
			tools.DebugPrintF("[INFO] Args: %s", strings.Join(args, " "))

			// db.Open()

			gHandler := gin.New()

			// RedirectFixedPath if enabled, the router tries to fix the current request path, if no
			// handle is registered for it.
			// First superfluous path elements like ../ or // are removed.
			// Afterwards the router does a case-insensitive lookup of the cleaned path.
			// If a handle can be found for this route, the router makes a redirection
			// to the corrected path with status code 301 for GET requests and 307 for
			// all other request methods.
			// For example /FOO and /..//Foo could be redirected to /foo.
			// RedirectTrailingSlash is independent of this option.
			gHandler.RedirectFixedPath = true

			// RemoveExtraSlash a parameter can be parsed from the URL even with extra slashes.
			// See the PR #1817 and issue #1644
			gHandler.RemoveExtraSlash = true

			gHandler.Use(middlewares.Configs(conf))
			gHandler.Use(middlewares.LoggerWithFormatter())
			gHandler.Use(middlewares.BasicAuth())
			gHandler.Use(middlewares.Cors())
			gHandler.Use(middlewares.Referer())
			gHandler.Use(middlewares.RateLimiter())
			gHandler.Use(middlewares.Gzip())
			gHandler.Use(middlewares.Header())
			gHandler.Use(gin.Recovery())

			if app.WwwRoot == "" {
				app.WwwRoot = "."
			}

			// gHandler.StaticFS("/", gin.Dir(r.WwwRoot, true))
			ghttp.StaticFS(&gHandler.RouterGroup, "/", http.Dir(app.WwwRoot))

			// HTTP SERVER
			host := app.Host
			port := app.Port

			var addrHTTP string
			// var addrHTTPS string

			if port != "" {
				addrHTTP = host + ":" + port
			} else {
				// Get the free port from 8080
				addrHTTP = ":" + gnet.GetAvailablePort(8080)
			}

			httpServer := &http.Server{
				Addr:           addrHTTP,
				Handler:        gHandler,
				ReadTimeout:    time.Duration(app.ReadTimeout) * time.Second,
				WriteTimeout:   time.Duration(app.WriteTimeout) * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			if !app.EnableHTTPS {
				tools.DebugPrintF("[INFO] HTTPS was disabled.")
				gracefulStart(httpServer)
				return
			}

			// HTTP SERVER
			httpsport := app.HTTPSPort

			var addrHTTPS string

			if httpsport != "" {
				addrHTTPS = host + ":" + httpsport
			} else {
				// Get the free port from 8443
				addrHTTPS = ":" + gnet.GetAvailablePort(8443)
			}

			// Load cert and key
			var certPEMBlock, keyPEMBlock []byte
			var err error
			var cert tls.Certificate
			var getCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)

			// load From Local Cert
			if app.HTTPSCertFile != "" && app.HTTPSKeyFile != "" {
				cert, err = tls.LoadX509KeyPair(app.HTTPSCertFile, app.HTTPSKeyFile)
			} else {
				err = errors.New("app.HTTPSCertFile is Empty or app.HTTPSKeyFile is empty")
			}

			// load From Auto Cert
			if err != nil {
				if app.Port == "80" && app.HTTPSPort == "443" && app.HTTPSDomains != nil && len(app.HTTPSDomains) > 0 && app.ContactEmail != "" {
					httpscertsdir := "certs"

					if app.HTTPSCertsDir != "" {
						httpscertsdir = app.HTTPSCertsDir
					}

					certManager := autocert.Manager{
						Prompt:     autocert.AcceptTOS,
						HostPolicy: autocert.HostWhitelist(app.HTTPSDomains...), //your domain here
						Cache:      autocert.DirCache(httpscertsdir),            //folder for storing certificates
						Email:      app.ContactEmail,
					}

					getCertificate = certManager.GetCertificate
				}
			}

			// load From Embde Cert
			if err != nil && getCertificate == nil {
				certPEMBlock, err = ghttps.GetTLSCerts().ReadFile("certs/server.pem")

				if err != nil {
					log.Fatal(err)
				}

				keyPEMBlock, err = ghttps.GetTLSCerts().ReadFile("certs/server.key")

				if err != nil {
					log.Fatal(err)
				}

				cert, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)

				if err != nil {
					log.Fatal(err)
				}
			}

			// Construct a tls.config
			tlsConfig := &tls.Config{
				Certificates:   []tls.Certificate{cert},
				GetCertificate: getCertificate,
				MinVersion:     tls.VersionTLS12,
				MaxVersion:     tls.VersionTLS13,
			}

			httpsServer := &http.Server{
				Addr:           addrHTTPS,
				TLSConfig:      tlsConfig,
				Handler:        gHandler,
				ReadTimeout:    time.Duration(app.ReadTimeout) * time.Second,
				WriteTimeout:   time.Duration(app.WriteTimeout) * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			gracefulStart(httpServer, httpsServer)
		},
	}

	rootCmd.Flags().StringVarP(configs.ConfigFile(), "config", "c", "", `If it is not set, we will try with development.(json/env/ini/yaml/toml/hcl/properties`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.Port, "port", "P", configs.GetConfigs().App.Port, `Port optionally specifies the TCP Port for the server to listen on,
in the form "host:port". If empty, ":port" (port 8080) is used.`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.Host, "host", "H", configs.GetConfigs().App.Host, `Host optionally specifies the Http Address for the server to listen on,
in the form "host:port". If empty, "host:" (host localhost) is used.`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.WwwRoot, "wwwroot", "w", configs.GetConfigs().App.WwwRoot, `By default, the wwwroot folder is treated as a web root folder. 
Static files can be stored in any folder under the web root and accessed with a relative path to that root.`)

	rootCmd.Flags().Int64VarP(&configs.GetConfigs().App.ReadTimeout, "read-timeout", "", configs.GetConfigs().App.ReadTimeout, `ReadTimeout is the maximum duration for reading the entire
request, including the body. A zero or negative value means
there will be no timeout.

Because ReadTimeout does not let Handlers make per-request
decisions on each request body's acceptable deadline or
upload rate, most users will prefer to use
ReadHeaderTimeout. It is valid to use them both.`)

	rootCmd.Flags().Int64VarP(&configs.GetConfigs().App.WriteTimeout, "write-timeout", "", configs.GetConfigs().App.WriteTimeout, `WriteTimeout is the maximum duration before timing out
writes of the response. It is reset whenever a new
request's header is read. Like ReadTimeout, it does not
let Handlers make decisions on a per-request basis.
A zero or negative value means there will be no timeout.`)

	rootCmd.Flags().BoolVarP(&configs.GetConfigs().App.EnableHTTPS, "enable-https", "", configs.GetConfigs().App.EnableHTTPS, `If it is set, we will enable https.`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.HTTPSPort, "https-port", "", configs.GetConfigs().App.HTTPSPort, `HTTPS PORT`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.HTTPSCertFile, "https-cert-file", "", configs.GetConfigs().App.HTTPSCertFile, `HTTPS Cert File`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.HTTPSKeyFile, "https-key-file", "", configs.GetConfigs().App.HTTPSKeyFile, `HTTPS Key File`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.HTTPSCertsDir, "https-cert-dir", "", configs.GetConfigs().App.HTTPSCertsDir, `HTTPS Cert Directory`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.ContactEmail, "contact-email", "", configs.GetConfigs().App.ContactEmail, `HTTPS Contact Email`)

	rootCmd.Flags().StringArrayVarP(&configs.GetConfigs().App.HTTPSDomains, "https-domains", "", configs.GetConfigs().App.HTTPSDomains, `HTTPS Domains`)

	rootCmd.Flags().BoolVarP(&configs.GetConfigs().App.Gzip, "gzip", "g", configs.GetConfigs().App.Gzip, `If it is set, we will compress with gzip.`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.AutoIndexTimeFormat, "autoindex-time-format", "", configs.GetConfigs().App.AutoIndexTimeFormat, `this is the AutoIndex Time Format.`)

	rootCmd.Flags().BoolVarP(&configs.GetConfigs().App.AutoIndexExactSize, "autoindex-exact-size", "", configs.GetConfigs().App.AutoIndexExactSize, `For the HTML format, specifies whether exact file sizes should be output in the directory listing,
or rather rounded to kilobytes, megabytes, and gigabytes.`)

	rootCmd.Flags().BoolVarP(&configs.GetConfigs().App.PreviewHTML, "preview-html", "", configs.GetConfigs().App.PreviewHTML, `For static web files, Whether preview them.`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.LogDir, "log-dir", "", configs.GetConfigs().App.LogDir, `The Log Directory which store access.log, error.log etc.`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.RateLimiter, "rate-limiter", "", configs.GetConfigs().App.RateLimiter, `Define a limit rate to several requests per hour.
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
	`)

	rootCmd.Flags().Int64VarP(&configs.GetConfigs().App.SpeedLimiter, "speed-limiter", "", configs.GetConfigs().App.SpeedLimiter, ` Specify  the  maximum  transfer  rate you want curl to use - for
downloads.
The given speed is measured in bytes/second, `)

	rootCmd.Flags().BoolVarP(&configs.GetConfigs().App.RefererLimiter, "referer-limiter", "", configs.GetConfigs().App.RefererLimiter, `Limit by referer`)

	rootCmd.Flags().BoolVarP(&configs.GetConfigs().App.Basic, "basic", "", configs.GetConfigs().App.Basic, `If it is set, we will use HTTP Basic authentication.

Used together with -u, --user <user:password>.

Providing --basic multiple times has no extra effect.

Example:`+env.ProjectName+`-u name:password --basic https://example.com`)

	rootCmd.Flags().StringVarP(&configs.GetConfigs().App.User, "user", "u", configs.GetConfigs().App.User, `Specify the user name and password to use for server authentication. 

The user name and passwords are split up  on  the  first  colon,
which  makes  it impossible to use a colon in the user name with
this option. The password can, still.`)
}

// Execute start the web server
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func welcome() {
	myFigure := figure.NewColorFigure("SnowdreamTech Static HTTP Server", "larry3d", "green", true)
	myFigure.Blink(1000, 1000, 0)
	figureString := myFigure.ColorString()

	fmt.Fprint(tools.DefaultGinWriter, figureString)
	fmt.Fprint(tools.DefaultGinWriter, "\n\n\n")
}

func gracefulStart(servers ...*http.Server) {
	var err error

	for _, server := range servers {
		// Initializing the server in a goroutine so that
		// it won't block the graceful shutdown handling below
		go func(server *http.Server) {
			if server.TLSConfig != nil {
				//https
				tools.DebugPrintF("[INFO] Listening and Serving HTTPS on %s\n", server.Addr)

				app := configs.GetAppConfig()

				if !app.EnableHTTPS {
					tools.DebugPrintF("[INFO] Hit CTRL-C to stop the server")
				}

				if err := server.ListenAndServeTLS(app.HTTPSCertFile, app.HTTPSKeyFile); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			} else {
				//http
				tools.DebugPrintF("[INFO] Listening and Serving HTTP on %s\n", server.Addr)

				app := configs.GetAppConfig()

				if app.EnableHTTPS {
					tools.DebugPrintF("[INFO] Hit CTRL-C to stop the server")
				}

				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}
		}(server)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	tools.DebugPrintF("[INFO] Shutting down servers...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, server := range servers {
		if err = server.Shutdown(ctx); err != nil {
			log.Fatal("The Web Server forced to shutdown: ", err)
		}
	}

	tools.DebugPrintF("[INFO] The Web Servers have been shut down.")
}
