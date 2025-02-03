package hud

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Manager interface {
	WritePacket(pk packet.Packet) error
}

func Hide(m Manager, e ...Element) {
	_ = m.WritePacket(&packet.SetHud{
		Elements:   toBytes(e...),
		Visibility: packet.HudVisibilityHide,
	})
}

func Reset(m Manager, e ...Element) {
	_ = m.WritePacket(&packet.SetHud{
		Elements:   toBytes(e...),
		Visibility: packet.HudVisibilityReset,
	})
}

func toBytes(e ...Element) []byte {
	var buff []byte

	for _, elem := range e {
		buff = append(buff, byte(elem))
	}

	return buff
}
