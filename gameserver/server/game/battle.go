package game

import "math"

type fightFunction func(uint32, uint32) uint32

func defaultF1Function(f1 uint32, f2 uint32) uint32 {
	return uint32(max(0, int32(f1)-int32(f2)))
}

func defaultF2Function(f1 uint32, f2 uint32) uint32 {
	return uint32(max(0, int32(f2)-int32(f1)))
}

func defaultTerrainF1Function(f1 uint32, f2 uint32) uint32 {
	return uint32(max(0, int32(f1)-(int32(f2)/2)))
}

func defaultTerrainF2Function(f1 uint32, f2 uint32) uint32 {
	return uint32(max(0, int32(f2)-2*int32(f1)))
}

func barrierTerrainF1Function(f1 uint32, f2 uint32) uint32 {
	if f1 == 0 {
		return 0
	}
	return uint32(max(0, int32(f1)-((int32(f2)-5)/2)))
}

func barrierTerrainF2Function(f1 uint32, f2 uint32) uint32 {
	if f1 == 0 {
		return f2
	}
	return uint32(max(0, int32(f2)-int32(f1)*2-5))
}

func battle2Group(g1 *SoldierGroup, g2 *SoldierGroup) {
	fight(&g1.count, &g2.count, defaultF1Function, defaultF2Function)
}

func fight(f1 *uint32, f2 *uint32, g1Function fightFunction, g2Function fightFunction) {
	var f1R uint32
	var f2R uint32
	f1R = g1Function(*f1, *f2)
	f2R = g2Function(*f1, *f2)
	*f1 = f1R
	*f2 = f2R
}

func battleNGroup(g1 *SoldierGroup, g2s []*SoldierGroup, g1Function fightFunction, g2Function fightFunction) *SoldierGroup {
	var combinedForce uint32 = 0
	var combinedForceMap map[*Player]uint32 = make(map[*Player]uint32)
	for _, g2 := range g2s {
		combinedForce += g2.count
		_, ok := combinedForceMap[g2.player]
		if ok {
			combinedForceMap[g2.player] += g2.count
		} else {
			combinedForceMap[g2.player] = g2.count
		}
	}
	var combinedForceCopy = combinedForce

	if g1 != nil {
		fight(&g1.count, &combinedForceCopy, g1Function, g2Function)
	}

	if combinedForceCopy == 0 {
		return g1
	}

	var combinedForceMapRemaining map[*Player]uint32 = make(map[*Player]uint32)

	for player, force := range combinedForceMap {
		var remaining uint32 = uint32(math.Round(float64(force) * float64(combinedForceCopy) / float64(combinedForce)))
		if remaining > 0 {
			combinedForceMapRemaining[player] = remaining
		}
	}

	for player, force := range combinedForceMapRemaining {
		var otherForces uint32 = 0
		for p2, f2 := range combinedForceMapRemaining {
			if p2 != player {
				otherForces += f2
			}
		}
		fight(&force, &otherForces, defaultF1Function, defaultF2Function)

		if force > 0 {
			return NewSoldierGroup(player, force)
		}
	}

	return g1
}
