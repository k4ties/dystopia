package rank

import "github.com/google/uuid"

var Player = def{}

type def struct{}

func (def) UUID() uuid.UUID       { return uuid.UUID{} }
func (def) Name() string          { return "" }
func (def) Format() string        { return "" }
func (def) DisplayRankName() bool { return false }
func (def) Priority() Priority    { return PriorityDefault }
