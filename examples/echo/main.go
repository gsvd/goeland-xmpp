package main

import (
	"github.com/gsvd/goeland-xmpp/address"
)

func main() {
	a, err := address.New(
		address.WithLocal("test"),
		address.WithDomain("gsvd.dev"),
		address.WithResource("client"),
	)
	if err != nil {
		panic(err)
	}
	println("Parsed address:", a.String())
	println("Bare address:", a.Bare().String())
	println("Local part:", a.Local().String())
}
