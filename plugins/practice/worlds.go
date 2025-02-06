package practice

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/instance/lobby"
	"github.com/k4ties/dystopia/plugins/practice/mw"
	"github.com/k4ties/dystopia/plugins/practice/user/hud"
)

var (
	lobbySpawn = mgl64.Vec3{-1, 49, -1}
	arenaSpawn = mgl64.Vec3{666, -34, 666}
)

func setupWorldsManager(m *plugin.Manager, path string) {
	m.Srv().World().SetSpawn(cube.Pos{100000, 1000, 100000})

	if err := mw.NewManager(m.Srv().World(), path, m.Logger()); err != nil {
		panic(err)
	}

	registerLobby(mw.M(), m)
	registerNodebuff(mw.M(), m)
}

func registerLobby(mn *mw.Manager, m *plugin.Manager) {
	w, ok := mn.World("lobby")
	if !ok {
		panic("must create a world called lobby")
	}

	if err := mn.SetSpawn("lobby", lobbySpawn); err != nil {
		panic(err)
	}

	hidden := []hud.Element{
		hud.Health, hud.Hunger,
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), lobby.Rotation(), hidden...)
	instance.Register("lobby", in)
}

func registerNodebuff(mn *mw.Manager, m *plugin.Manager) {
	w, ok := mn.World("arena")
	if !ok {
		panic("must create a world called arena")
	}

	if err := mn.SetSpawn("arena", arenaSpawn); err != nil {
		panic(err)
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), cube.Rotation{})
	instance.Register("nodebuff", in)
}
