// coalesce

// coalesce.go

package main

import (
	"runtime"

	"github.com/nytopop/utron"

	_ "git.echoesofthe.net/nytopop/coalesce/controllers"
	_ "git.echoesofthe.net/nytopop/coalesce/models"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	utron.Run()
}
