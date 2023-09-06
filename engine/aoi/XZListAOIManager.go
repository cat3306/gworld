package aoi

type xzAoi struct {
	aoi          *AOI
	neighbors    map[*xzAoi]struct{}
	xPrev, xNext *xzAoi
	yPrev, yNext *xzAoi
	markVal      int
}

// XZListAOIManager is an implementation of AOICalculator using XZ lists
type XZListAOIManager struct {
	xSweepList *xAOIList
	zSweepList *yAOIList
}

// NewXZListAOIManager creates a new XZListAOIManager
func NewXZListAOIManager(aoiDist float32) Manager {
	return &XZListAOIManager{
		xSweepList: newXAOIList(aoiDist),
		zSweepList: newYAOIList(aoiDist),
	}
}

// Enter is called when Entity enters Space
func (aoiman *XZListAOIManager) Enter(aoi *AOI, x, y float32) {
	xzaoi := &xzAoi{
		aoi:       aoi,
		neighbors: map[*xzAoi]struct{}{},
	}
	aoi.x, aoi.y = x, y
	aoi.implData = xzaoi
	aoiman.xSweepList.Insert(xzaoi)
	aoiman.zSweepList.Insert(xzaoi)
	aoiman.adjust(xzaoi)
}

// Leave is called when Entity leaves Space
func (aoiman *XZListAOIManager) Leave(aoi *AOI) {
	xzaoi := aoi.implData.(*xzAoi)
	aoiman.xSweepList.Remove(xzaoi)
	aoiman.zSweepList.Remove(xzaoi)
	aoiman.adjust(xzaoi)
}

// Moved is called when Entity moves in Space
func (aoiman *XZListAOIManager) Moved(aoi *AOI, x, y float32) {
	oldX := aoi.x
	oldY := aoi.y
	aoi.x, aoi.y = x, y
	xzaoi := aoi.implData.(*xzAoi)
	if oldX != x {
		aoiman.xSweepList.Move(xzaoi, oldX)
	}
	if oldY != y {
		aoiman.zSweepList.Move(xzaoi, oldY)
	}
	aoiman.adjust(xzaoi)
}

// adjust is called by Entity to adjust neighbors
func (aoiman *XZListAOIManager) adjust(aoi *xzAoi) {
	aoiman.xSweepList.Mark(aoi)
	aoiman.zSweepList.Mark(aoi)
	// AOI marked twice are neighbors
	for neighbor := range aoi.neighbors {
		if neighbor.markVal == 2 {
			// neighbors kept
			neighbor.markVal = -2 // mark this as neighbor
		} else { // markVal < 2
			// was neighbor, but not any more
			delete(aoi.neighbors, neighbor)
			aoi.aoi.callback.OnLeaveAOI(neighbor.aoi)
			delete(neighbor.neighbors, aoi)
			neighbor.aoi.callback.OnLeaveAOI(aoi.aoi)
		}
	}

	// travel in X list again to find all new neighbors, whose markVal == 2
	aoiman.xSweepList.GetClearMarkedNeighbors(aoi)
	// travel in Z list again to unmark all
	aoiman.zSweepList.ClearMark(aoi)
}
