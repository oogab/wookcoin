package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/oogab/wookcoin/explorer"
	"github.com/oogab/wookcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to wookcoin\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port: 	Set the PORT of the server\n")
	fmt.Printf("-mode: 	Choose between 'html' and 'rest'\n\n")
	// defer를 이행하고 나서 Exit를 해야한다.
	// os.Exit(0)
	// runtime.Goexit은 모든 함수를 제거하지만 그 전에 defer를 먼저 이행한다.
	runtime.Goexit()
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		// start rest api
		rest.Start(*port)
	case "html":
		// start html explorer
		explorer.Start(*port)
	default:
		usage()
	}
}
