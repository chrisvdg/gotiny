package main

import (
	"github.com/spf13/pflag"
)

func main() {
	var port int16
	var tlsPort int16

	pflag.Int16VarP(&port, "port", "p", 8080, "http listen port")
	pflag.Int16VarP(&tlsPort, "tlsport", "t", 8443, "https listen port")

}
