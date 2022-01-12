package pprof4gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http/pprof"
	"path"
	"strings"
)

func Run(prefix string, route gin.IRoutes) {
	run(prefix, route, "")
}

func RunAddress(prefix, address string) {
	run(prefix, nil, address)
}

func run(prefix string, route gin.IRoutes, address string) {
	var nEngine *gin.Engine
	if route == nil {
		nEngine = gin.New()
		nEngine.Use(gin.Logger(), gin.Recovery())
		route = nEngine
	}

	prefix = strings.TrimSpace(prefix)
	if prefix != "" {
		prefix = cleanPath(prefix)
	}

	var p = path.Join(prefix, "/debug/pprof/*action")

	route.Any(p, func(c *gin.Context) {
		var nPath = strings.TrimPrefix(c.Request.URL.Path, prefix)
		c.Request.URL.Path = nPath

		switch c.Request.URL.Path {
		case "/debug/pprof/":
			pprof.Index(c.Writer, c.Request)
		case "/debug/pprof/cmdline":
			pprof.Cmdline(c.Writer, c.Request)
		case "/debug/pprof/profile":
			pprof.Profile(c.Writer, c.Request)
		case "/debug/pprof/symbol":
			pprof.Symbol(c.Writer, c.Request)
		case "/debug/pprof/trace":
			pprof.Trace(c.Writer, c.Request)
		default:
			pprof.Index(c.Writer, c.Request)
		}
	})

	if nEngine != nil {
		if address == "" {
			address = "127.0.0.1:0"
		}

		var addr, err = net.ResolveTCPAddr("tcp", address)
		if err != nil {
			addr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:0")
		}

		if addr != nil && addr.IP == nil {
			addr.IP = net.IPv4(127, 0, 0, 1)
		}

		var listener, _ = net.ListenTCP("tcp", addr)
		var nPath = fmt.Sprintf("%s%s/debug/pprof", listener.Addr().String(), prefix)
		log.Printf("pprof is listening on http://%s", nPath)
		go nEngine.RunListener(listener)
	}
}

func cleanPath(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	//trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			//trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	//if trailing && w > 1 {
	//	bufApp(&buf, p, w, '/')
	//	w++
	//}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}
