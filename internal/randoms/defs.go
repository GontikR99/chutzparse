package randoms

type RollGroup struct {
	Min   int32
	Max   int32
	Rolls []*CharacterRoll
}

type CharacterRoll struct {
	Character string
	Age       string
	Value     int32
}

const ChannelChange = "RandomsChange"
