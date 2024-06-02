package game

type direction uint8

const (
	From direction = iota
	To
)

type groupDirection struct {
	Group     *SoldierGroup
	Direction direction
	progress  uint8
}

func (g *groupDirection) getTrueLength(totalLength uint8) uint8 {
	if g.Direction == To {
		return g.progress
	}
	return totalLength - g.progress
}

func (g *groupDirection) isOtherInFrontOfMe(other *groupDirection, totalLength uint8) bool {
	var actualPos1 uint8 = g.getTrueLength(totalLength)
	var actualPos2 uint8 = other.getTrueLength(totalLength)
	if actualPos1 == actualPos2 {
		return true
	}
	if g.Direction == To {
		return actualPos1+1 == actualPos2
	} else {
		return actualPos1-1 == actualPos2
	}
}

type Pipe struct {
	firstLocation  Node
	secondLocation Node
	groups         []*groupDirection
	length         uint8
	pendingGroups  [2]*SoldierGroup
}

func (p *Pipe) GetGroupCount() int {
	return len(p.groups)
}

func NewPipe(firstLocation Node, secondLocation Node, length uint8) *Pipe {
	return &Pipe{
		firstLocation:  firstLocation,
		secondLocation: secondLocation,
		length:         length + 1,
		pendingGroups:  [2]*SoldierGroup{},
		groups:         []*groupDirection{},
	}
}

func (p *Pipe) Tick() {
	var groupToRemove []int
	for i, g := range p.groups {
		g.progress++
		if g.progress == p.length {
			groupToRemove = append(groupToRemove, i)
		}
	}

	for i := len(groupToRemove); i > 0; i-- {
		var indexToRemove int = groupToRemove[i-1]
		var g = p.groups[indexToRemove]
		p.groups = append(p.groups[:indexToRemove], p.groups[indexToRemove+1:]...)
		if g.Direction == To {
			p.secondLocation.AddToActionQueue(g.Group, p)
		} else {
			p.firstLocation.AddToActionQueue(g.Group, p)
		}
	}

	if p.pendingGroups[0] != nil {
		pendingGroup := p.pendingGroups[0]
		p.pendingGroups[0] = nil
		p.groups = append(p.groups, &groupDirection{
			Group:     pendingGroup,
			Direction: To,
			progress:  1,
		})
	}

	if p.pendingGroups[1] != nil {
		pendingGroup := p.pendingGroups[1]
		p.pendingGroups[1] = nil
		p.groups = append(p.groups, &groupDirection{
			Group:     pendingGroup,
			Direction: From,
			progress:  1,
		})
	}
}

func (p *Pipe) PostTick() {
	for _, g1 := range p.groups {
		for _, g2 := range p.groups {
			if g1.Group.player != g2.Group.player && g1.Direction != g2.Direction && g1.Direction == To {
				if g1.isOtherInFrontOfMe(g2, p.length) {
					battle2Group(g1.Group, g2.Group)
				}
			}
		}
	}

	var groups []*groupDirection
	for _, g := range p.groups {
		if g.Group.count > 0 {
			groups = append(groups, g)
		}
	}
	p.groups = groups
}

func (p *Pipe) AddToActionQueue(soldierGroup *SoldierGroup, from Location) {
	if from == p.firstLocation {
		p.pendingGroups[0] = soldierGroup
	} else if from == p.secondLocation {
		p.pendingGroups[1] = soldierGroup
	}
}
