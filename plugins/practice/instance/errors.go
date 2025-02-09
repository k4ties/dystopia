package instance

import (
	"errors"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

const kickMessageFormat = "Asynchronous <red>error</red>\nCode: <grey>%s</grey>"

type ErrorCode error

func Kick(p *player.Player, e ErrorCode) {
	p.Disconnect(text.Colourf(kickMessageFormat, toTitle(e.Error())))
}

// ErrorSponge reason: fail to sync connection in instance (no connection on plugin manager)
var ErrorSponge = errors.New("sponge")

// ErrorPussy reason: fail to sync with database
var ErrorPussy = errors.New("pussy")

// ErrorAngus reason: fail to sync user data (tick event)
var ErrorAngus = errors.New("angus")

// toTitle uppercases the first letter of the string
func toTitle(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
