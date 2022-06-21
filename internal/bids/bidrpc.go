package bids

//go:generate ../../build/rpcgen bidrpc.go

// Bids provides the renderer the ability to interact with the bid tracking subsytem
type Bids interface {
	RefreshDKP() (int32, error)
	AuctionActive() (bool, error)
	HasGuildDump() (bool, error)
	FetchBids() ([]*ItemBids, error)

	Start() error
	End() error
}
