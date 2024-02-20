/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ansible-rekey/cmd"
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
