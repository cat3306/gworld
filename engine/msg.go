package engine

type TryConnectMsg struct {
	NetWork string
	Addr    string
}

func CheckInnerMsg(i *InnerMsg) bool {
	return len(i.ClientIds) != 0
}
