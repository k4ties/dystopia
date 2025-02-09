package user

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/google/uuid"
	"github.com/k4ties/dystopia/plugins/practice/rank"
	"maps"
	"slices"
	"sync"
	"time"
)

type User struct {
	d *data
}

func New(c session.Conn, o OfflineUser) *User {
	var os, input int

	if c == nil {
		os = 1
		input = 1
	}

	return &User{d: newData(OS(os), InputMode(input), o)}
}

func (u *User) Dystopia() *data {
	return u.d
}

func newData(os OS, i InputMode, o OfflineUser) *data {
	d := &data{}

	d.SwitchInputMode(i)
	d.switchOS(os)

	loadFromOffline(d, o)
	return d
}

type data struct {
	name string
	xuid string

	uuid uuid.UUID

	ips  []string
	dids []string

	online atomic.Bool

	os    atomic.Value[OS]
	input atomic.Value[InputMode]

	kills, deaths atomic.Int64
	firstJoin     time.Time

	rank atomic.Value[rank.Rank]
}

func (d *data) Online() bool {
	return d.online.Load()
}

func (d *data) ApplyFromOnlinePlayer(p *player.Player, c session.Conn) {
	d.online.Store(true)

	if !slices.Contains(d.ips, p.Addr().String()) {
		d.AddIP(p.Addr().String())
	}
	if !slices.Contains(d.dids, p.DeviceID()) {
		d.AddDeviceID(p.DeviceID())
	}

	d.switchOS(OS(c.ClientData().DeviceOS))
	d.SwitchInputMode(InputMode(c.ClientData().CurrentInputMode))
}

func (d *data) SyncQuit() {
	d.online.Store(false)
	d.os.Store(-1)
	d.input.Store(65535)
}

func loadFromOffline(d *data, o OfflineUser) {
	d.ips = o.IPS
	d.dids = o.DeviceIDS

	d.kills.Store(int64(o.Kills))
	d.deaths.Store(int64(o.Deaths))

	d.firstJoin = o.FirstJoin
	d.rank.Store(o.Rank)
}

func (d *data) Rank() rank.Rank {
	return d.rank.Load()
}

func (d *data) SetRank(r rank.Rank) {
	d.rank.Store(r)
}

func (d *data) AddIP(ip string) {
	d.ips = append(d.ips, ip)
}

func (d *data) AddDeviceID(did string) {
	d.dids = append(d.dids, did)
}

func (d *data) DeviceIDs() []string {
	return d.dids
}

func (d *data) switchOS(new OS) {
	if d.Online() {
		d.os.Store(new)
	}
}

func (d *data) SwitchInputMode(new InputMode) {
	if d.Online() {
		d.input.Store(new)
	}
}

func (d *data) IPS() []string {
	return d.ips
}

func (d *data) Kill() {
	if d.Online() {
		d.kills.Inc()
	}
}

func (d *data) Death() {
	if d.Online() {
		d.deaths.Inc()
	}
}

func (d *data) Kills() int64 {
	return d.kills.Load()
}

func (d *data) Deaths() int64 {
	return d.deaths.Load()
}

func (d *data) FirstJoin() time.Time {
	return d.firstJoin
}

func (d *data) OS() OS {
	if !d.Online() {
		return -1
	}

	return d.os.Load()
}

func (d *data) InputMode() InputMode {
	if !d.Online() {
		return InputMode(65535)
	}

	return d.input.Load()
}

func ValidOS(o OS) bool {
	return o != -1
}

func ValidInputMode(m InputMode) bool {
	return m != 65535
}

func (d *data) offline() OfflineUser {
	return OfflineUser{
		IPS:       d.ips,
		FirstJoin: d.firstJoin,
		Name:      d.name,
		XUID:      d.xuid,
		UUID:      d.uuid.String(),
		Deaths:    int(d.deaths.Load()),
		Kills:     int(d.kills.Load()),
		Rank:      d.Rank(),
	}
}

var pool = struct {
	v  map[uuid.UUID]*User
	mu sync.RWMutex
}{
	v: make(map[uuid.UUID]*User),
}

func Register(u *User) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.v[u.d.uuid] = u
}

func List() []*User {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return slices.Collect(maps.Values(pool.v))
}

func ByName(n string) (*User, bool) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	for _, u := range List() {
		if u.d.name == n {
			return u, true
		}
	}

	return nil, false
}

func ByXUID(xuid string) (*User, bool) {
	for _, u := range List() {
		if u.d.xuid == xuid {
			return u, true
		}
	}

	return nil, false
}

func ByDeviceID(id string) (*User, bool) {
	for _, u := range List() {
		if slices.Contains(u.d.DeviceIDs(), id) {
			return u, true
		}
	}

	return nil, false
}

func ByUUID(uuid uuid.UUID) (*User, bool) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	u, ok := pool.v[uuid]
	return u, ok
}

type XUID string
type Name string
type DeviceID string

func Lookup(v any) (*User, bool) {
	switch v := v.(type) {
	case uuid.UUID:
		return ByUUID(v)
	case Name:
		return ByName(string(v))
	case XUID:
		return ByXUID(string(v))
	case DeviceID:
		return ByDeviceID(string(v))
	}

	return nil, false
}

//func RegisterFromDB() {
//	d := database.D()
//
//	for _, acc := range d.Accounts() {
//		New(nil, OfflineUserFromAccount(acc))
//	}
//}
