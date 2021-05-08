package main

import (
	"github.com/alicek106/go-ec2-ssh-autoconnect/pkg/cmd"
	"log"
)

func main() {
	err := cmd.NewCommand().Execute()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
