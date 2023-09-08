package gameobject

import "github.com/cat3306/goworld/engine/aoi"

type ISpace interface {
	OnInit()    // Called when initializing space struct, override to initialize custom space fields
	OnCreate()  // Called when space is created
	OnDestroy() // Called just before space is destroyed
	OnGameObjectEnter(obj IObject)
	OnGameObjectLeave(ob IObject)
	OnGameReady()
}

type Space struct {
	aoiMgr      aoi.Manager
	gameObjects GameObjectSet
}

func (s *Space) OnInit() {

}
func (s *Space) OnCreated() {

}
func (s *Space) OnDestroy() {

}
func (s *Space) OnGameObjectEnter(obj IObject) {

}

func (s *Space) OnGameObjectLeave(obj IObject) {

}
func (s *Space) OnGameReady() {

}