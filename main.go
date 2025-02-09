package main

import (
	"github.com/k4ties/dystopia/dystopia"
	"log/slog"
)

func main() {
	l := slog.Default()
	conf := dystopia.MustReadConfig("config.yaml")

	d := dystopia.New(l, conf)
	d.Start()
}
