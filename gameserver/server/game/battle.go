package game

func battle2Group(g1 *SoldierGroup, g2 *SoldierGroup) {
	var g1Count uint32 = uint32(max(int64(g1.count)-int64(g2.count), 0))
	var g2Count uint32 = uint32(max(int64(g2.count)-int64(g1.count), 0))
	g1.count = g1Count
	g2.count = g2Count
}
