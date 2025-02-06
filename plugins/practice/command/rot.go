package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Rot struct {
	onlyPlayer
	// no arguments
}

func (Rot) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	rot := p(s).Rotation()

	o.Printf(text.Colourf("Pitch: <grey>%f</grey>\nYaw: <grey>%f</grey>", rot.Pitch(), rot.Yaw()))
}

func init() {
	cmd.Register(cmd.New("rot", "", nil, Rot{}))
}
