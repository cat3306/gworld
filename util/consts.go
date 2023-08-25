package util

type ClusterType int

var (
	CallGate   = MethodHash("CallClient")
	CallClient = MethodHash("CallClient")
)

const (
	ClusterTypeGate       ClusterType = 0
	ClusterTypeDispatcher ClusterType = 1
	ClusterTypeGame       ClusterType = 2
)
const (
	ChanPacketSize = 10000
)
const (
	DispatcherConnMgrKey = "DispatcherConnMgrKey"
	GameConnMgrKey       = "GameConnMgrKey"
)
const (
	GateClientMgrKey      = "GateClientMgrKey"
	GameClientProxyMgrKey = "GameClientProxyMgrKey"
)

func (c ClusterType) String() string {
	if c == ClusterTypeGate {
		return "gate"
	} else if c == ClusterTypeDispatcher {
		return "dispatcher"
	} else if c == ClusterTypeGame {
		return "game"
	}
	return "invalid cluster type"
}

const (
	MethodSetDispatcherType = "SetDispatcherType"
)