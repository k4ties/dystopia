package rank

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var Owner = owner{}

type owner struct{}

func (owner) UUID() uuid.UUID       { return uuid.MustParse("5a74f657-9663-403b-b74c-eadecec332e2") }
func (owner) Name() string          { return "Owner" }
func (owner) Format() string        { return text.Bold + text.Red }
func (owner) DisplayRankName() bool { return true }
func (owner) Priority() Priority    { return PriorityOwner }
