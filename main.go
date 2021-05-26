package main

import (
	"bytes"
	"fmt"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
)

func main() {
	c := collection.Collection{}

	err := c.AddFile("test.txt", bytes.NewBufferString("test"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Create("support-collector.zip")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.WriteZIP(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file.Close()
}
