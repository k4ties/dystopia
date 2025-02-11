package practice

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/dystopia/embeddable"
	"github.com/k4ties/dystopia/plugins/practice/ffa"
	"github.com/k4ties/dystopia/plugins/practice/ffa/fist"
	"github.com/k4ties/dystopia/plugins/practice/ffa/gapple"
	"github.com/k4ties/dystopia/plugins/practice/ffa/nodebuff"
	"github.com/k4ties/dystopia/plugins/practice/ffa/sumo"
	"github.com/k4ties/dystopia/plugins/practice/instance"
	"github.com/k4ties/dystopia/plugins/practice/mw"
	"github.com/k4ties/dystopia/plugins/practice/user/hud"
	"time"
)

type WorldConfig struct {
	Config struct {
		Name string

		Rotation [2]float64
		Spawn    [3]float64

		HeightThreshold instance.HeightThresholdConfig
	}
}

var (
	//go:embed lobby.json
	lobbyJson []byte
	//go:embed nodebuff.json
	arenaJson []byte
	//go:embed gapple.json
	gappleJson []byte
	//go:embed sumo.json
	sumoJson []byte
	//go:embed fist.json
	fistJson []byte
)

var (
	lobbyConfig  = embeddable.MustJSON[WorldConfig](lobbyJson)
	arenaConfig  = embeddable.MustJSON[WorldConfig](arenaJson)
	gappleConfig = embeddable.MustJSON[WorldConfig](gappleJson)
	sumoConfig   = embeddable.MustJSON[WorldConfig](sumoJson)
	fistConfig   = embeddable.MustJSON[WorldConfig](fistJson)
)

func setupWorldsManager(m *plugin.Manager, path string) {
	if err := mw.NewManager(m.Srv().World(), path, m.Logger()); err != nil {
		panic(err)
	}

	registerLobby(mw.M(), m)
	// ffa
	registerNodebuff(mw.M(), m)
	registerGapple(mw.M(), m)
	registerSumo(mw.M(), m)
	registerFist(mw.M(), m)
}

func registerLobby(mn *mw.Manager, m *plugin.Manager) {
	name := lobbyConfig.Config.Name

	w, ok := mn.World(name)
	if !ok {
		panic("no world with specified name on the config")
	}

	if err := mn.SetSpawn(name, lobbyConfig.Config.Spawn); err != nil {
		panic(err)
	}

	hidden := []hud.Element{
		hud.Health, hud.Hunger,
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), lobbyConfig.Config.Rotation, lobbyConfig.Config.HeightThreshold, hidden...)
	instance.Register(name, in)
}

func registerNodebuff(mn *mw.Manager, m *plugin.Manager) {
	name := arenaConfig.Config.Name

	w, ok := mn.World(name)
	if !ok {
		panic("no world with specified name on the config")
	}

	if err := mn.SetSpawn(name, arenaConfig.Config.Spawn); err != nil {
		panic(err)
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), arenaConfig.Config.Rotation, arenaConfig.Config.HeightThreshold)

	f := ffa.New(in.(*instance.Impl), nodebuff.Kit, ffa.Config{
		Name: "NoDebuff",
		Icon: "textures/items/potion_bottle_splash_heal.png",

		PearlCooldown: time.Second * 16,
	})

	instance.Register("nodebuff", f)
}

func registerGapple(mn *mw.Manager, m *plugin.Manager) {
	name := gappleConfig.Config.Name

	w, ok := mn.World(name)
	if !ok {
		panic("no world with specified name on the config")
	}

	if err := mn.SetSpawn(name, gappleConfig.Config.Spawn); err != nil {
		panic(err)
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), gappleConfig.Config.Rotation, gappleConfig.Config.HeightThreshold)

	f := ffa.New(in.(*instance.Impl), gapple.Kit, ffa.Config{
		Name: "GApple",
		Icon: "textures/items/apple_golden.png",
	})

	instance.Register("gapple", f)
}

func registerSumo(mn *mw.Manager, m *plugin.Manager) {
	name := sumoConfig.Config.Name

	w, ok := mn.World(name)
	if !ok {
		panic("no world with specified name on the config")
	}

	if err := mn.SetSpawn(name, sumoConfig.Config.Spawn); err != nil {
		panic(err)
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), sumoConfig.Config.Rotation, sumoConfig.Config.HeightThreshold)

	f := ffa.New(in.(*instance.Impl), sumo.Kit, ffa.Config{
		Name: "Sumo",
		Icon: "textures/items/slimeball.png",
	})

	instance.Register("sumo", f)
}

func registerFist(mn *mw.Manager, m *plugin.Manager) {
	name := fistConfig.Config.Name

	w, ok := mn.World(name)
	if !ok {
		panic("no world with specified name on the config")
	}

	if err := mn.SetSpawn(name, fistConfig.Config.Spawn); err != nil {
		panic(err)
	}

	in := instance.New(w, world.GameModeSurvival, m.Logger(), fistConfig.Config.Rotation, fistConfig.Config.HeightThreshold)

	f := ffa.New(in.(*instance.Impl), fist.Kit, ffa.Config{
		Name: "Fist",
		Icon: "textures/items/beef_cooked.png",
	})

	instance.Register("fist", f)
}
