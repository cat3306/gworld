package util

type ClusterType int

var (
	CallClient = MethodHash("CallClient")
	//CallClient        = MethodHash("CallClient")
	SetClientProperty = MethodHash("SetClientProperty")
	ClientOnConnect   = MethodHash("OnConnect")
)

const (
	ClusterTypeGate ClusterType = 0
	ClusterTypeGame ClusterType = 1
)
const (
	ChanPacketSize = 10000
)
const (
	GateClientMgrKey      = "GateClientMgrKey"
)

func (c ClusterType) String() string {
	if c == ClusterTypeGate {
		return "gate"
	} else if c == ClusterTypeGame {
		return "game"
	}
	return "invalid cluster type"
}

const (
	ClientAuth = "ClientAuth"
)
const (
	RoomId   = "RoomId"
	ClientId = "ClientId"
)

const (
	ContextInnerMsgKey = "InnerMsg"
)
