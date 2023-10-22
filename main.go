package main

import (
	"log"
)

func main() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(listEngineCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
