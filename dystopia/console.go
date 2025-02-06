package dystopia

import (
	"bufio"
	"github.com/df-mc/dragonfly/server"
	"github.com/k4ties/dystopia/plugins/practice/handlers/whitelist"
	"log/slog"
	"os"
	"strings"
)

type WhitelistConfig struct {
	Whitelist struct {
		Enabled bool
		Players []string
	}
}

type ConsoleConfig struct {
	Logger     *slog.Logger
	ConfigPath string

	Server *server.Server
}

func StartConsole(c ConsoleConfig) {
	scanner := bufio.NewScanner(os.Stdin)

	l := c.Logger

	go func() {
		for scanner.Scan() {
			args := strings.Split(scanner.Text(), " ")

			switch args[0] {
			case "wl", "whitelist", "w":
				if len(args) <= 1 {
					l.Info("Usage: wl add|delete|toggle")
					continue
				}

				switch args[1] {
				case "add", "+":
					name := strings.Join(args[2:], " ")
					if name == "" {
						l.Error("Please specify a name")
						continue
					}

					whitelist.Add(name)
					l.Info("Successfully added player " + name + " to whitelist.")

					cf := MustReadConfig(c.ConfigPath).AddWhitelistPlayer(name)
					cf.syncWithFile(c.ConfigPath)
				case "delete", "-":
					name := strings.Join(args[2:], " ")
					if name == "" {
						l.Error("Please specify a name")
						continue
					}

					whitelist.Remove(name)
					l.Info("Successfully deleted player " + name + " from whitelist.")

					cf := MustReadConfig(c.ConfigPath).RemoveWhitelistPlayer(name)
					cf.syncWithFile(c.ConfigPath)
				case "toggle":
					whitelist.Toggle()
					l.Info("Successfully toggled whitelist.")

					cf := MustReadConfig(c.ConfigPath).ToggleWhitelist()
					cf.syncWithFile(c.ConfigPath)
				default:
					l.Error("Unknown argument: " + args[1])
				}
			case "stop":
				os.Exit(2)
			default:
				l.Error("Unknown command: " + args[0])
			}
		}
	}()
}
