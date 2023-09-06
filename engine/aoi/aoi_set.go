package aoi

type Set map[*AOI]struct{}

func (s Set) Add(a *AOI) {
	s[a] = struct{}{}
}

func (s Set) Remove(a *AOI) {
	delete(s, a)
}

func (s Set) Contains(a *AOI) (ok bool) {
	_, ok = s[a]
	return
}
