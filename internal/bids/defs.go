package bids

type CharacterStat struct {
	Rank       string
	Balance    int32
	Attendance []string
}

type ItemBid struct {
	CalculatedBid int32
	BidMessages   []string
}

type AnnotatedBid struct {
	Character string
	Bid       ItemBid
	Stat      CharacterStat
}

type ItemBids struct {
	Item string
	Bids []*AnnotatedBid
}

type MemberInfo struct {
	Rank    string
	Comment string
}

const UnspecifiedItem = "Unspecified"
const ChannelChange = "BidsChange"
