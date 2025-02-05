// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package http

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/snowdreamtech/gserver/pkg/env"
)

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//
//	router.Static("/static", "/var/www")
func Static(group *gin.RouterGroup, relativePath, root string) gin.IRoutes {
	return group.StaticFS(relativePath, gin.Dir(root, false))
}

// StaticFS works just like `Static()` but a custom `FileSystem` can be used instead.
// Gin by default uses: gin.Dir()
func StaticFS(group *gin.RouterGroup, relativePath string, fs http.FileSystem) gin.IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving ad static folder")
	}
	handler := createStaticHandler(group, relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return nil
}

func createStaticHandler(group *gin.RouterGroup, relativePath string, fs http.FileSystem) gin.HandlerFunc {
	absolutePath := calculateAbsolutePath(group, relativePath)
	fileServer := http.StripPrefix(absolutePath, FileServer(fs))

	return func(c *gin.Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)

			c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

			w := c.Writer
			title := "404 Not Found"

			fmt.Fprintf(w, "<html>\n")

			fmt.Fprintf(w, "<head>\n")
			fmt.Fprintf(w, "<title>\n")
			fmt.Fprintf(w, "%s\n", title)
			fmt.Fprintf(w, "</title>\n")
			fmt.Fprintf(w, "</head>\n")

			fmt.Fprintf(w, "<style>\n")
			fmt.Fprintf(w, "body %s\n", "{display: flex;min-height: 100vh;flex-direction: column; margin:0px; padding:0px 8px;}")
			fmt.Fprintf(w, "hr %s\n", "{display:block;border: 0;width:100%;height: 1px;background-color:#555555;clear:both;}")
			fmt.Fprintf(w, "span %s\n", "{display: inline-block;width:300px;}")
			fmt.Fprintf(w, ".link %s\n", "{text-decoration: none;color: #000; padding:0 5px}")
			fmt.Fprintf(w, ".header %s\n", "{display: flex;flex-direction: column;flex: 0 0 auto;}")
			fmt.Fprintf(w, ".content %s\n", "{display: flex;flex-direction: column;flex: 1 0 auto;}")
			fmt.Fprintf(w, ".footer %s\n", "{display: flex; justify-content: center;align-items: center;flex-direction: row; flex: 0 0 auto; padding-bottom:10px;'}")
			fmt.Fprintf(w, "</style>\n")

			fmt.Fprintf(w, "<body>\n")
			fmt.Fprintf(w, "<div class=\"header\">\n")
			fmt.Fprintf(w, "<center>\n")
			fmt.Fprintf(w, "<h1>\n")
			fmt.Fprintf(w, "%s\n", title)
			fmt.Fprintf(w, "</h1>\n")
			fmt.Fprintf(w, "</center>\n")
			fmt.Fprintf(w, "</div>\n")
			fmt.Fprintf(w, "<div class=\"content\">\n")
			fmt.Fprintf(w, "<hr>\n")
			fmt.Fprintf(w, "<center>\n")

			fmt.Fprintf(w, "<center>\n")
			fmt.Fprintf(w, "%s/%s\n", env.ProjectName, env.GitTag)
			fmt.Fprintf(w, "</center>\n")
			fmt.Fprintf(w, "</div>\n")
			fmt.Fprintf(w, "<div class=\"footer\">\n")
			fmt.Fprintf(w, "Powered by")
			fmt.Fprintf(w, "<a href=\"%s\" class=\"link\">%s</a>\n", "https://github.com/snowdreamtech/gserver", "SnowdreamTech Static HTTP Server")
			fmt.Fprintf(w, "</div>\n")
			fmt.Fprintf(w, "</body>\n")
			fmt.Fprintf(w, "</html>\n")
			// c.handlers = group.engine.noRoute
			// // Reset index
			// c.index = -1
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func calculateAbsolutePath(group *gin.RouterGroup, relativePath string) string {
	return joinPaths(group.BasePath(), relativePath)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
