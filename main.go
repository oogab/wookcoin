package main

import (
	"github.com/oogab/wookcoin/explorer"
	"github.com/oogab/wookcoin/rest"
)

func main() {
	// go explorer.Start(3000)
	explorer.Start(5500)
	rest.Start(4000)
}
