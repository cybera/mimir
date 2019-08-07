package main

import (
	"log"

	"github.com/cybera/mimir/cmd"
)

func main() {
	log.SetFlags(0)
	cmd.Execute()
}
