package domain

type AccountType uint8

func (a AccountType) AsUint8() uint8 {
	return uint8(a)
}

const (
	AccountTypeUnknown = iota
	AccountTypeReward
	AccountTypeSystem
)

type CreditItem struct {
	Uid         int64
	Account     int64
	AccountType AccountType
	Amt         int64
	Currency    string
}

type Credit struct {
	Biz   string
	BizId int64
	Items []CreditItem
}
