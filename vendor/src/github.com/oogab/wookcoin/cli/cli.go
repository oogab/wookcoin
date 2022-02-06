package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/oogab/wookcoin/explorer"
	"github.com/oogab/wookcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 우크 코인\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port: Set the PORT of the server\n")
	fmt.Printf("-mode: Choose between 'html' and 'rest'\n")
	os.Exit(1)
}

func Start() {

	// if len(os.Args) < 2 {
	// 	usage()
	// }

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
		explorer.Start(*port)
	default:
		usage()
	}

	fmt.Println(*port, *mode)

	// fmt.Println(os.Args[2:])
	// rest := flag.NewFlagSet("rest", flag.ExitOnError)

	// portFlag := rest.Int("port", 4000, "Sets the port of the server")

	// switch os.Args[1] {
	// case "explorer":
	// 	fmt.Println("Start Explorer")
	// case "rest":
	// 	rest.Parse(os.Args[2:])
	// default:
	// 	usage()
	// }

	// if rest.Parsed() {
	// 	fmt.Println(portFlag)
	// 	fmt.Println("Start server")
	// }
}