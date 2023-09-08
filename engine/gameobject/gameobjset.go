package gameobject

type GameObjectSet map[string]IObject

func (g GameObjectSet) Add(obj IObject) {
	g[obj.GetId()] = obj
}
func (g GameObjectSet) Get(id string) IObject {
	return g[id]
}
