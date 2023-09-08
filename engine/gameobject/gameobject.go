package gameobject

type IObject interface {
	OnCreated(id string)
	// Migration
	OnMigrateOut()
	OnMigrateIn()
	// Freeze && Restore
	OnFreeze()
	OnRestored()
	// Space Operations
	OnEnterSpace()
	OnLeaveSpace() //space *Space
	// Client Notifications
	OnClientConnected()
	OnClientDisconnected() // Called when Client disconnected

	//DescribeEntityType() // Define entity attributes in this function
	OnMove(pos Vector3, rot Vector3)
	GetId() string
}

type GameObject struct {
	Id           string
	Tag          string
	destroyed    bool
	Position     Vector3
	Rotation     Vector3
	InterestedIn GameObjectSet
	InterestedBy GameObjectSet
	Components   GameObjectSet
	client       *GameClient
	Property     map[string]interface{}
}

func (g *GameObject) OnMigrateOut() {

}
func (g *GameObject) OnMigrateIn() {

}
func (g *GameObject) OnClientDisconnected() {

}
func (g *GameObject) OnClientConnected() {

}
func (g *GameObject) OnCreated(id string) {
	g.Components = map[string]IObject{}
	g.InterestedIn = map[string]IObject{}
	g.InterestedBy = map[string]IObject{}
	g.Property = map[string]interface{}{}
	g.Id = id
}
func (g *GameObject) GetId() string {
	return g.Id
}

func (g *GameObject) OnMove(pos Vector3, rot Vector3) {
	g.Position = pos
	g.Rotation = rot
}
func (g *GameObject) OnRestored() {

}
func (g *GameObject) OnFreeze() {

}
func (g *GameObject) OnEnterSpace() {

}
func (g *GameObject) OnLeaveSpace() {

}
