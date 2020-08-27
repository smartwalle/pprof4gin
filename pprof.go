package pprof4gin

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

func Run(r gin.IRoutes) {
	r.GET("/debug/pprof/*action", func(c *gin.Context) {
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
}
