package aoi

import "log"

type TowerAOIManager struct {
	minX, maxX, minY, maxY float32
	towerRange             float32
	towers                 [][]tower
	xTowerNum, yTowerNum   int
}

func (aoiMan *TowerAOIManager) Enter(aoi *AOI, x, y float32) {
	aoi.x, aoi.y = x, y
	obj := &aoiObj{aoi: aoi}
	aoi.implData = obj

	aoiMan.visitWatchedTowers(x, y, aoi.dist, func(tower *tower) {
		tower.addWatcher(obj)
	})

	t := aoiMan.getTowerXY(x, y)
	t.addObj(obj, nil)
}

func (aoiMan *TowerAOIManager) Leave(aoi *AOI) {
	obj := aoi.implData.(*aoiObj)
	obj.tower.removeObj(obj, true)

	aoiMan.visitWatchedTowers(aoi.x, aoi.y, aoi.dist, func(tower *tower) {
		tower.removeWatcher(obj)
	})
}

func (aoiMan *TowerAOIManager) Moved(aoi *AOI, x, y float32) {
	oldx, oldy := aoi.x, aoi.y
	aoi.x, aoi.y = x, y
	obj := aoi.implData.(*aoiObj)
	t0 := obj.tower
	t1 := aoiMan.getTowerXY(x, y)

	if t0 != t1 {
		t0.removeObj(obj, false)
		t1.addObj(obj, t0)
	}

	oximin, oximax, oyimin, oyimax := aoiMan.getWatchedTowers(oldx, oldy, aoi.dist)
	ximin, ximax, yimin, yimax := aoiMan.getWatchedTowers(x, y, aoi.dist)

	for xi := oximin; xi <= oximax; xi++ {
		for yi := oyimin; yi <= oyimax; yi++ {
			if xi >= ximin && xi <= ximax && yi >= yimin && yi <= yimax {
				continue
			}

			tower := &aoiMan.towers[xi][yi]
			tower.removeWatcher(obj)
		}
	}

	for xi := ximin; xi <= ximax; xi++ {
		for yi := yimin; yi <= yimax; yi++ {
			if xi >= oximin && xi <= oximax && yi >= oyimin && yi <= oyimax {
				continue
			}

			tower := &aoiMan.towers[xi][yi]
			tower.addWatcher(obj)
		}
	}
}

func (aoiMan *TowerAOIManager) transXY(x, y float32) (int, int) {
	xi := int((x - aoiMan.minX) / aoiMan.towerRange)
	yi := int((y - aoiMan.minY) / aoiMan.towerRange)
	return aoiMan.normalizeXi(xi), aoiMan.normalizeYi(yi)
}

func (aoiMan *TowerAOIManager) normalizeXi(xi int) int {
	if xi < 0 {
		xi = 0
	} else if xi >= aoiMan.xTowerNum {
		xi = aoiMan.xTowerNum - 1
	}
	return xi
}

func (aoiMan *TowerAOIManager) normalizeYi(yi int) int {
	if yi < 0 {
		yi = 0
	} else if yi >= aoiMan.yTowerNum {
		yi = aoiMan.yTowerNum - 1
	}
	return yi
}

func (aoiMan *TowerAOIManager) getTowerXY(x, y float32) *tower {
	xi, yi := aoiMan.transXY(x, y)
	return &aoiMan.towers[xi][yi]
}

func (aoiMan *TowerAOIManager) getWatchedTowers(x, y float32, aoiDistance float32) (int, int, int, int) {
	ximin, yimin := aoiMan.transXY(x-aoiDistance, y-aoiDistance)
	ximax, yimax := aoiMan.transXY(x+aoiDistance, y+aoiDistance)
	//aoiTowerNum := int(aoiDistance/aoiMan.towerRange) + 1
	//ximid, yimid := aoiMan.transXY(x, y)
	//ximin, ximax := aoiMan.normalizeXi(ximid-aoiTowerNum), aoiMan.normalizeXi(ximid+aoiTowerNum)
	//yimin, yimax := aoiMan.normalizeYi(yimid-aoiTowerNum), aoiMan.normalizeYi(yimid+aoiTowerNum)
	return ximin, ximax, yimin, yimax
}

func (aoiMan *TowerAOIManager) visitWatchedTowers(x, y float32, aoiDistance float32, f func(*tower)) {
	ximin, ximax, yimin, yimax := aoiMan.getWatchedTowers(x, y, aoiDistance)
	for xi := ximin; xi <= ximax; xi++ {
		for yi := yimin; yi <= yimax; yi++ {
			tower := &aoiMan.towers[xi][yi]
			f(tower)
		}
	}
}

func (aoiMan *TowerAOIManager) init() {
	numXSlots := int((aoiMan.maxX-aoiMan.minX)/aoiMan.towerRange) + 1
	aoiMan.xTowerNum = numXSlots
	numYSlots := int((aoiMan.maxY-aoiMan.minY)/aoiMan.towerRange) + 1
	aoiMan.yTowerNum = numYSlots
	aoiMan.towers = make([][]tower, numXSlots)
	for i := 0; i < numXSlots; i++ {
		aoiMan.towers[i] = make([]tower, numYSlots)
		for j := 0; j < numYSlots; j++ {
			aoiMan.towers[i][j].init()
		}
	}
}

func NewTowerAOIManager(minX, maxX, minY, maxY float32, towerRange float32) Manager {
	aoiMan := &TowerAOIManager{minX: minX, maxX: maxX, minY: minY, maxY: maxY, towerRange: towerRange}
	aoiMan.init()

	return aoiMan
}

type tower struct {
	objs     map[*aoiObj]struct{}
	watchers map[*aoiObj]struct{}
}

func (t *tower) init() {
	t.objs = map[*aoiObj]struct{}{}
	t.watchers = map[*aoiObj]struct{}{}
}

func (t *tower) addObj(obj *aoiObj, fromOtherTower *tower) {
	obj.tower = t
	t.objs[obj] = struct{}{}
	if fromOtherTower == nil {
		for watcher := range t.watchers {
			if watcher == obj {
				continue
			}
			watcher.aoi.callback.OnEnterAOI(obj.aoi)
		}
	} else {
		// obj moved from other tower to this tower
		for watcher := range fromOtherTower.watchers {
			if watcher == obj {
				continue
			}
			if _, ok := t.watchers[watcher]; ok {
				continue
			}
			watcher.aoi.callback.OnLeaveAOI(obj.aoi)
		}
		for watcher := range t.watchers {
			if watcher == obj {
				continue
			}
			if _, ok := fromOtherTower.watchers[watcher]; ok {
				continue
			}
			watcher.aoi.callback.OnEnterAOI(obj.aoi)
		}
	}
}

func (t *tower) removeObj(obj *aoiObj, notifyWatchers bool) {
	obj.tower = nil
	delete(t.objs, obj)
	if notifyWatchers {
		for watcher := range t.watchers {
			if watcher == obj {
				continue
			}
			watcher.aoi.callback.OnLeaveAOI(obj.aoi)
		}
	}
}

func (t *tower) addWatcher(obj *aoiObj) {
	if _, ok := t.watchers[obj]; ok {
		log.Panicf("duplicate add watcher")
	}
	t.watchers[obj] = struct{}{}
	// now obj can see all objs under this tower
	for neighbor := range t.objs {
		if neighbor == obj {
			continue
		}
		obj.aoi.callback.OnEnterAOI(neighbor.aoi)
	}
}

func (t *tower) removeWatcher(obj *aoiObj) {
	if _, ok := t.watchers[obj]; !ok {
		log.Panicf("duplicate remove watcher")
	}

	delete(t.watchers, obj)
	for neighbor := range t.objs {
		if neighbor == obj {
			continue
		}
		obj.aoi.callback.OnLeaveAOI(neighbor.aoi)
	}
}

type aoiObj struct {
	aoi   *AOI
	tower *tower
}
