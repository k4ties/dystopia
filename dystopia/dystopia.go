package dystopia

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player/chat"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/df-plugin/example/npc"
	"github.com/k4ties/dystopia/plugins/practice"
	"github.com/k4ties/dystopia/plugins/practice/handlers"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"log/slog"
	"os"
	"path/filepath"
)

type Dystopia struct {
	l *slog.Logger
	c Config

	m *plugin.Manager
}

func New(l *slog.Logger, c Config) *Dystopia {
	d := &Dystopia{c: c, l: l}
	d.setup()

	return d
}

func (d *Dystopia) setup() {
	d.m = plugin.NewManager(plugin.ManagerConfig{
		Logger:     d.l,
		UserConfig: d.c.convert(),
		SubName:    "",
		Packs:      d.loadPacks(),
	})

	d.m = d.m.WithPlayerProvider(practice.Provider(cube.Rotation{180, 0}, d.m))

	d.m.ToggleStatusCommand()
	d.m.Register(practice.Plugin(d.loginHandler(), d.c.Advanced.CachePath+"/worlds"), npc.Plugin())
}

func (d *Dystopia) loginHandler() *handlers.Login {
	return practice.LoginHandler(d.c.Whitelist.Enabled, d.c.Whitelist.Players...).(*handlers.Login)
}

func (d *Dystopia) loadPacks() (pool []*resource.Pack) {
	path := d.c.Resources.Path

	dir, err := os.ReadDir(path)
	if err != nil {
		panic("couldn't read dir while loading packs: " + err.Error())
	}

	for i, f := range dir {
		pathTo := filepath.Join(path, f.Name())
		pack, err := resource.ReadPath(pathTo)
		if err != nil {
			d.l.Error("dystopia: resource packs: cannot load directory: " + pathTo)
			continue
		}

		pool = append(pool, pack.WithContentKey(d.c.Resources.ContentKey))
		d.l.Debug(fmt.Sprintf("dystopia: loaded %d/%d packs", i, len(dir)))
	}

	return
}

func (d *Dystopia) Start() {
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	d.m.Srv().CloseOnProgramEnd()
	d.m.ListenServer()
}
