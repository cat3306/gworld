package gameobject

import "fmt"

type GameObjectSet map[string]IObject

func (g GameObjectSet) Add(obj IObject) {
	if obj == nil {
		fmt.Println("asd")
	}
	if g == nil {
		fmt.Println("ad")
	}
	g[obj.GetId()] = obj
}
func (g GameObjectSet) Get(id string) IObject {
	return g[id]
}
