package main

import (
	"github.com/smartwalle/pprof4gin"
)

func main() {
	pprof4gin.Run("test", nil)
	select {}
}
