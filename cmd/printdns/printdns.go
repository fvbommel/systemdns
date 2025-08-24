package main

import (
	"fmt"
	"log"

	"github.com/fvbommel/systemdns"
)

func main() {
	servers, err := systemdns.GetSystemDNS()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("System DNS resolvers:")
	for _, s := range servers {
		fmt.Printf(" - %s\n", s)
	}
}
