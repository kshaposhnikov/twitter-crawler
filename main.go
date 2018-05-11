package main

import (
	"runtime"
	"github.com/kshaposhnikov/twitter-crawler/command"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	command.Execute()
}
