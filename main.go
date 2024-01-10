/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ansible-rekey/cmd"
	"log"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
