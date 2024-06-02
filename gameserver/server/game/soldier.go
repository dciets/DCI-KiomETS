package game

type SoldierGroup struct {
	player *Player
	count  uint32
}

func NewSoldierGroup(player *Player, count uint32) *SoldierGroup {
	return &SoldierGroup{player: player, count: count}
}

func (sg *SoldierGroup) Count() uint32 {
	return sg.count
}
