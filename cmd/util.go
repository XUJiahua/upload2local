package cmd

import (
	"fmt"
	"os"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
