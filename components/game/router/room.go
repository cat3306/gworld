package router

type ClientInfo struct {
	ConnId string
}
type Room struct {
	maxNum    int    //人数
	pwd       string //密码
	joinState bool   //是否能加入
	gameState bool   //游戏状态
	scene     int    //游戏场景
	Id        string
	clients   map[string]*ClientInfo
}
