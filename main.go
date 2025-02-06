package main

import (
	"fmt"
	"github.com/k4ties/dystopia/dystopia"
	"log/slog"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic recovered: %v\n", r)
		}
	}()

	l := slog.Default()
	conf := dystopia.MustReadConfig("config.yaml")

	d := dystopia.New(l, conf)
	d.Start()
}
