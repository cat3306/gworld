package util

type ClusterType int

var (
	CallGate   = MethodHash("CallClient")
	CallClient = MethodHash("CallClient")
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
	GameClientProxyMgrKey = "GameClientProxyMgrKey"
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
	ClientAuth = "auth"
)
