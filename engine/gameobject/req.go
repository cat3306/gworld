package gameobject

type PlayerPos struct {
	Vector3
	NetObjId string  `json:"NetObjId"`
	Yaw      float32 `json:"Yaw"`
	CX       float32 `json:"CX"`
}
