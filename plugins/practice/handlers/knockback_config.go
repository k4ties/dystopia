package handlers

import (
	_ "embed"
	"github.com/k4ties/dystopia/dystopia/embeddable"
)

type kbConfig struct {
	KnockBack struct {
		Force    float64
		Height   float64
		Immunity int
	}
}

//go:embed knockback.json
var knockbackConfig []byte

var KnockbackConfig = embeddable.MustJSON[kbConfig](knockbackConfig)
