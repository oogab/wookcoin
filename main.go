package main

import (
	"github.com/oogab/wookcoin/cli"
	"github.com/oogab/wookcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
