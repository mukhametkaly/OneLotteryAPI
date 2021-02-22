package Lottery

import (
	"math/rand"
	"time"
)

func (l Lottery) Chek ()  bool{
	return !((l.Prize == "") || (l.TextMessage == "") || (l.LotName == "") || &l.Raffler.UserID == nil)
}

func (l *Lottery) Play () bool {
	size := len(l.PlayerIDs)
	if size == 0 {
		return false
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(size)
	l.Winner = &l.PlayerIDs[index]
	return true

}