package events

type ReadEvent struct {
	Uid int64
	Aid int64
}

type ReadEventV1 struct {
	Uids []int64
	Aids []int64
}
