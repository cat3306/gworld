package aoi

type yAOIList struct {
	aoiDist float32
	head    *xzAoi
	tail    *xzAoi
}

func newYAOIList(aoiDist float32) *yAOIList {
	return &yAOIList{aoiDist: aoiDist}
}

func (sl *yAOIList) Insert(aoi *xzAoi) {
	insertCoord := aoi.aoi.y
	if sl.head != nil {
		p := sl.head
		for p != nil && p.aoi.y < insertCoord {
			p = p.yNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert xzAoi at the end of list
			tail := sl.tail
			tail.yNext = aoi
			aoi.yPrev = tail
			sl.tail = aoi
		} else { // otherwise, p >= xzAoi, insert xzAoi before p
			prev := p.yPrev
			aoi.yNext = p
			p.yPrev = aoi
			aoi.yPrev = prev

			if prev != nil {
				prev.yNext = aoi
			} else { // p is the head, so xzAoi should be the new head
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}

func (sl *yAOIList) Remove(aoi *xzAoi) {
	prev := aoi.yPrev
	next := aoi.yNext
	if prev != nil {
		prev.yNext = next
		aoi.yPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.yPrev = prev
		aoi.yNext = nil
	} else {
		sl.tail = prev
	}
}

func (sl *yAOIList) Move(aoi *xzAoi, oldCoord float32) {
	coord := aoi.aoi.y
	if coord > oldCoord {
		// moving to next ...
		next := aoi.yNext
		if next == nil || next.aoi.y >= coord {
			// no need to adjust in list
			return
		}
		prev := aoi.yPrev
		//fmt.Println(1, prev, next, prev == nil || prev.yNext == xzAoi)
		if prev != nil {
			prev.yNext = next // remove xzAoi from list
		} else {
			sl.head = next // xzAoi is the head, trim it
		}
		next.yPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		prev, next = next, next.yNext
		for next != nil && next.aoi.y < coord {
			prev, next = next, next.yNext
			//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.X < coord && (next == nil || next.X >= coord), so insert between prev and next
		prev.yNext = aoi
		aoi.yPrev = prev
		if next != nil {
			next.yPrev = aoi
		} else {
			sl.tail = aoi
		}
		aoi.yNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := aoi.yPrev
		if prev == nil || prev.aoi.y <= coord {
			// no need to adjust in list
			return
		}

		next := aoi.yNext
		if next != nil {
			next.yPrev = prev
		} else {
			sl.tail = prev // xzAoi is the head, trim it
		}
		prev.yNext = next // remove xzAoi from list

		next, prev = prev, prev.yPrev
		for prev != nil && prev.aoi.y > coord {
			next, prev = prev, prev.yPrev
		}
		// no we have next.X > coord && (prev == nil || prev.X <= coord), so insert between prev and next
		next.yPrev = aoi
		aoi.yNext = next
		if prev != nil {
			prev.yNext = aoi
		} else {
			sl.head = aoi
		}
		aoi.yPrev = prev
	}
}

func (sl *yAOIList) Mark(aoi *xzAoi) {
	prev := aoi.yPrev
	coord := aoi.aoi.y

	minCoord := coord - sl.aoiDist
	for prev != nil && prev.aoi.y >= minCoord {
		prev.markVal += 1
		prev = prev.yPrev
	}

	next := aoi.yNext
	maxCoord := coord + sl.aoiDist
	for next != nil && next.aoi.y <= maxCoord {
		next.markVal += 1
		next = next.yNext
	}
}

func (sl *yAOIList) ClearMark(aoi *xzAoi) {
	prev := aoi.yPrev
	coord := aoi.aoi.y

	minCoord := coord - sl.aoiDist
	for prev != nil && prev.aoi.y >= minCoord {
		prev.markVal = 0
		prev = prev.yPrev
	}

	next := aoi.yNext
	maxCoord := coord + sl.aoiDist
	for next != nil && next.aoi.y <= maxCoord {
		next.markVal = 0
		next = next.yNext
	}
}
