package dystopia

import (
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/dystopia/plugins/practice"
	"github.com/k4ties/dystopia/plugins/practice/handlers"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"log/slog"
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
	d.m = plugin.NewManager(d.l, text.Colourf(d.c.Dystopia.SubName), d.c.convert())
	d.m.ToggleStatusCommand()

	d.m.Register(practice.Plugin(d.loginHandler(), d.c.Advanced.CachePath+"/worlds"))
}

func (d *Dystopia) loginHandler() *handlers.Login {
	return practice.LoginHandler(d.c.Whitelist.Enabled, d.c.Whitelist.Players...).(*handlers.Login)
}

func (d *Dystopia) Start() {
	d.m.Srv().CloseOnProgramEnd()
	d.m.ListenServer()
}
