package serialisation

type GroupSerialisation struct {
	OwnerIndex   int  `json:"ownerIndex"`
	SoldierCount uint `json:"soldierCount"`
	Length       uint `json:"length"`
	Upward       bool `json:"upward"`
}

func NewGroupSerialisation(ownerIndex int, soldierCount uint, length uint, upward bool) *GroupSerialisation {
	return &GroupSerialisation{
		OwnerIndex:   ownerIndex,
		SoldierCount: soldierCount,
		Length:       length,
		Upward:       upward,
	}
}

type PipeSerialisation struct {
	Length   uint                 `json:"length"`
	First    uint                 `json:"first"`
	Second   uint                 `json:"second"`
	Soldiers []GroupSerialisation `json:"soldiers"`
}

func NewPipeSerialisation(length uint, first uint, second uint, groupLength int) *PipeSerialisation {
	return &PipeSerialisation{
		Length:   length,
		First:    first,
		Second:   second,
		Soldiers: make([]GroupSerialisation, groupLength),
	}
}
