package main

import (
	"github.com/gsvd/goeland-xmpp/address"
)

func main() {
	a, err := address.Parse("test@gsvd.dev/client")
	if err != nil {
		panic(err)
	}
	println("Parsed address:", a.String())
	println("Bare address:", a.Bare().String())
	println("Local part:", a.Local().String())
}
