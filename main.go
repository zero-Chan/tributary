package main

import (
	"fmt"

	"tributary/config"
)

func main() {
	ag, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	ag.Start()

	return
}
