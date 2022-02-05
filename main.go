package main

import (
	explorer "github.com/oogab/wookchain/explorer/templates"
	"github.com/oogab/wookchain/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
