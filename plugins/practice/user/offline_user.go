package user

import (
	"github.com/k4ties/dystopia/plugins/practice/rank"
	"time"
)

type OfflineUser struct {
	IPS       []string
	DeviceIDS []string

	FirstJoin time.Time

	Name       string
	XUID, UUID string

	Deaths, Kills int
	Rank          rank.Rank
}

//func OfflineUserFromAccount(a database.Account) OfflineUser {
//	return OfflineUser{
//		IPS:       a.IPs(),
//		DeviceIDS: a.DeviceIDs(),
//		FirstJoin: time.Unix(a.FirstJoin, 0),
//		Name:      a.Name,
//		XUID:      a.XUID,
//		UUID:      a.UUID,
//		Deaths:    a.Deaths,
//		Kills:     a.Kills,
//		Rank:      a.Rank(),
//	}
//}
