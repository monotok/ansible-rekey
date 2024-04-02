/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/monotok/ansible-utils/cmd"
	"log"
)

func main() {
	rekeyCli := cmd.NewRekeyCli()
	cmd := cmd.NewCliCommand(rekeyCli)
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
