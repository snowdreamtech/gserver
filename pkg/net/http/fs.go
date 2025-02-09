// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// HTTP file system request handler

package http

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"html"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/juju/ratelimit"
	"github.com/snowdreamtech/gserver/pkg/configs"
	"github.com/snowdreamtech/gserver/pkg/io"
)

// condResult is the result of an HTTP request precondition check.
// See https://tools.ietf.org/html/rfc7232 section 3.
type condResult int

const (
	condNone condResult = iota
	condTrue
	condFalse
)

type fileHandler struct {
	root http.FileSystem
}

type anyDirs interface {
	len() int
	name(i int) string
	isDir(i int) bool
	size(i int) int64
	modtime(i int) time.Time
}

type fileInfoDirs []fs.FileInfo

func (d fileInfoDirs) len() int                { return len(d) }
func (d fileInfoDirs) isDir(i int) bool        { return d[i].IsDir() }
func (d fileInfoDirs) name(i int) string       { return d[i].Name() }
func (d fileInfoDirs) size(i int) int64        { return d[i].Size() }
func (d fileInfoDirs) modtime(i int) time.Time { return d[i].ModTime() }

type dirEntryDirs []fs.DirEntry

func (d dirEntryDirs) len() int          { return len(d) }
func (d dirEntryDirs) isDir(i int) bool  { return d[i].IsDir() }
func (d dirEntryDirs) name(i int) string { return d[i].Name() }
func (d dirEntryDirs) size(i int) int64 {
	if d.isDir(i) {
		return -1
	}

	fileinfo, err := d[i].Info()

	if err != nil {
		return 0
	}

	return fileinfo.Size()
}

func (d dirEntryDirs) modtime(i int) time.Time {
	fileinfo, err := d[i].Info()
	if err != nil {
		return time.Time{}
	}

	return fileinfo.ModTime()
}

// FileServer returns a handler that serves HTTP requests
// with the contents of the file system rooted at root.
//
// As a special case, the returned file server redirects any request
// ending in "/index.html" to the same path, without the final
// "index.html".
//
// To use the operating system's file system implementation,
// use http.Dir:
//
//	http.Handle("/", http.FileServer(http.Dir("/tmp")))
//
// To use an fs.FS implementation, use http.FS to convert it:
//
//	http.Handle("/", http.FileServer(http.FS(fsys)))
func FileServer(root http.FileSystem) http.Handler {
	return &fileHandler{root}
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const options = http.MethodOptions + ", " + http.MethodGet + ", " + http.MethodHead

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		if !strings.HasPrefix(r.URL.Path, "/") {
			r.URL.Path = "/" + r.URL.Path
		}
		serveFile(w, r, f.root, path.Clean(r.URL.Path), true)

	case http.MethodOptions:
		w.Header().Set("Allow", options)

	default:
		w.Header().Set("Allow", options)
		http.Error(w, "read-only", http.StatusMethodNotAllowed)
	}
}

// name is '/'-separated, not filepath.Separator.
func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string, redirect bool) {
	indexPage := "/index.html"

	app := configs.GetAppConfig()

	if !app.PreviewHTML {
		indexPage = "/index.html123456789"
	}

	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(r.URL.Path, indexPage) {
		localRedirect(w, r, "./")
		return
	}

	f, err := fs.Open(filepath.Clean(name))
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if redirect {
		// redirect to canonical path: / at end of directory url
		// r.URL.Path always begins with /
		url := r.URL.Path
		if d.IsDir() {
			if url[len(url)-1] != '/' {
				localRedirect(w, r, path.Base(url)+"/")
				return
			}
		} else {
			if url[len(url)-1] == '/' {
				localRedirect(w, r, "../"+path.Base(url))
				return
			}
		}
	}

	if d.IsDir() {
		url := r.URL.Path
		// redirect if the directory name doesn't end in a slash
		if url == "" || url[len(url)-1] != '/' {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}

		// use contents of index.html for directory, if present
		index := strings.TrimSuffix(name, "/") + indexPage
		ff, err := fs.Open(filepath.Clean(index))
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				d = dd
				f = ff
			}
		}
	}

	// Still a directory? (we didn't find an index.html file)
	if d.IsDir() {
		if checkIfModifiedSince(r, d.ModTime()) == condFalse {
			writeNotModified(w)
			return
		}
		setLastModified(w, d.ModTime())
		dirList(w, r, f)
		return
	}

	if app.SpeedLimiter <= 0 {
		http.ServeContent(w, r, d.Name(), d.ModTime(), f)
	} else {
		bucket := ratelimit.NewBucketWithRate(float64(app.SpeedLimiter), app.SpeedLimiter)
		readseeker := io.ReadSeeker(f, f, bucket)
		http.ServeContent(w, r, d.Name(), d.ModTime(), readseeker)
	}
}

func dirList(w http.ResponseWriter, r *http.Request, f http.File) {
	// Prefer to use ReadDir instead of Readdir,
	// because the former doesn't require calling
	// Stat on every entry of a directory on Unix.
	var dirs anyDirs
	var err error
	if d, ok := f.(fs.ReadDirFile); ok {
		var list dirEntryDirs
		list, err = d.ReadDir(-1)
		dirs = list
	} else {
		var list fileInfoDirs
		list, err = f.Readdir(-1)
		dirs = list
	}

	if err != nil {
		//logf(r, "http: error reading directory: %v", err)
		http.Error(w, "http.Error reading directory", http.StatusInternalServerError)
		return
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs.name(i) < dirs.name(j) })

	app := configs.GetAppConfig()
	timeformat := "2006-01-02 15:04:05"

	if app.AutoIndexTimeFormat != "" {
		timeformat = app.AutoIndexTimeFormat
	}

	title := fmt.Sprintf("Index of %s", html.EscapeString(r.URL.String()))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, "<html>\n")

	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<title>\n")
	fmt.Fprintf(w, "%s\n", title)
	fmt.Fprintf(w, "</title>\n")
	fmt.Fprintf(w, "</head>\n")

	fmt.Fprintf(w, "<style>\n")
	fmt.Fprintf(w, "body %s\n", "{display: flex;min-height: 100vh;flex-direction: column; margin:0px; padding:0px 8px;}")
	fmt.Fprintf(w, "hr %s\n", "{display:block;border: 0;width:100%;height: 1px;background-color:#555555;clear:both;}")
	fmt.Fprintf(w, ".link %s\n", "{text-decoration: none;color: #000; padding:0 5px}")
	fmt.Fprintf(w, ".item %s\n", "{display: flex; flex-direction: row;justify-content: flex-start;align-items: flex-start;}")
	fmt.Fprintf(w, ".item-file %s\n", "{display: flex; flex-grow: 6; min-width:600px;max-width:600px; overflow:hidden;text-overflow:ellipsis;white-space:nowrap}")
	fmt.Fprintf(w, ".item-time %s\n", "{display: flex; flex-grow: 3; min-width:300px;max-width:300px;}")
	fmt.Fprintf(w, ".item-size %s\n", "{display: flex; flex-grow: 2; min-width:200px;max-width:200px;}")
	fmt.Fprintf(w, ".header %s\n", "{display: flex;flex-direction: column; justify-content: flex-start;align-items: flex-start;flex: 0 0 auto;}")
	fmt.Fprintf(w, ".content %s\n", "{display: flex;flex-direction: column; justify-content: flex-start;align-items: flex-start;flex: 1 0 auto;}")
	fmt.Fprintf(w, ".footer %s\n", "{display: flex; justify-content: center;align-items: center;flex-direction: row; flex: 0 0 auto; padding-bottom:10px;'}")
	fmt.Fprintf(w, "</style>\n")

	fmt.Fprintf(w, "<body>\n")
	fmt.Fprintf(w, "<div class=\"header\">\n")
	fmt.Fprintf(w, "<h1>\n")
	fmt.Fprintf(w, "%s\n", title)
	fmt.Fprintf(w, "</h1>\n")
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "<div class=\"content\">\n")
	fmt.Fprintf(w, "<hr>\n")

	// fmt.Fprintf(w, "<pre>\n")

	fmt.Fprintf(w, "<item class=\"item\">\n")
	fmt.Fprintf(w, "<span class=\"item-file\">")
	fmt.Fprintf(w, "<a href=\"%s\">../</a>\n", "../")
	fmt.Fprintf(w, "</span>\n")
	fmt.Fprintf(w, "</item>\n")

	for i, n := 0, dirs.len(); i < n; i++ {
		fmt.Fprintf(w, "<item class=\"item\">\n")

		name := dirs.name(i)
		if dirs.isDir(i) {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		if dirs.isDir(i) {
			fmt.Fprintf(w, "<span class=\"item-file\"><a href=\"%s\">%s</a></span><span class=\"item-time\">%s</span><span class=\"item-size\">%s</span>\n", url.String(), htmlReplacer.Replace(name), dirs.modtime(i).Format(timeformat), "-")
		} else {
			if app.AutoIndexExactSize {
				fmt.Fprintf(w, "<span class=\"item-file\"><a href=\"%s\">%s</a></span><span class=\"item-time\">%s</span><span class=\"item-size\">%d</span>\n", url.String(), htmlReplacer.Replace(name), dirs.modtime(i).Format(timeformat), dirs.size(i))
			} else {
				fmt.Fprintf(w, "<span class=\"item-file\"><a href=\"%s\">%s</a></span><span class=\"item-time\">%s</span><span class=\"item-size\">%s</span>\n", url.String(), htmlReplacer.Replace(name), dirs.modtime(i).Format(timeformat), units.HumanSize(float64(dirs.size(i))))
			}
		}

		fmt.Fprintf(w, "</item>\n")
	}
	// fmt.Fprintf(w, "</pre>\n")
	fmt.Fprintf(w, "<hr>\n")
	fmt.Fprintf(w, "</div>\n")
	fmt.Fprintf(w, "<div class=\"footer\">\n")
	fmt.Fprintf(w, "Powered by")
	fmt.Fprintf(w, "<a href=\"%s\" class=\"link\">%s</a>\n", "https://github.com/snowdreamtech/gserver", "SnowdreamTech Static HTTP Server")
	fmt.Fprintf(w, "</div>\n")
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

// toHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that toHTTPError does not
// actually return err.http.Error(), since msg and http.Status are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func toHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server http.Error", http.StatusInternalServerError
}

func checkIfModifiedSince(r *http.Request, modtime time.Time) condResult {
	if r.Method != "GET" && r.Method != "HEAD" {
		return condNone
	}
	ims := r.Header.Get("If-Modified-Since")
	if ims == "" || isZeroTime(modtime) {
		return condNone
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return condNone
	}
	// The Last-Modified header truncates sub-second precision so
	// the modtime needs to be truncated too.
	modtime = modtime.Truncate(time.Second)
	if ret := modtime.Compare(t); ret <= 0 {
		return condFalse
	}
	return condTrue
}

func writeNotModified(w http.ResponseWriter) {
	// RFC 7232 section 4.1:
	// a sender SHOULD NOT generate representation metadata other than the
	// above listed fields unless said metadata exists for the purpose of
	// guiding cache updates (e.g., Last-Modified might be useful if the
	// response does not have an ETag field).
	h := w.Header()
	delete(h, "Content-Type")
	delete(h, "Content-Length")
	delete(h, "Content-Encoding")
	if h.Get("Etag") != "" {
		delete(h, "Last-Modified")
	}
	w.WriteHeader(http.StatusNotModified)
}

// TimeFormat is the time format to use when generating times in HTTP
// headers. It is like time.RFC1123 but hard-codes GMT as the time
// zone. The time being formatted must be in UTC for Format to
// generate the correct format.
//
// For parsing this time format, see ParseTime.
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

var unixEpochTime = time.Unix(0, 0)

// isZeroTime reports whether t is obviously unspecified (either zero or Unix()=0).
func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

func setLastModified(w http.ResponseWriter, modtime time.Time) {
	if !isZeroTime(modtime) {
		w.Header().Set("Last-Modified", modtime.UTC().Format(TimeFormat))
	}
}

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)
