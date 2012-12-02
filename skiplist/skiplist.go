package skiplist

// XXX: Fix dice!
// XXX: Dump dot file

// Ref: William Pugh "Skip lists: a probabilistic alternative to balanced trees"

func (s *SkipList) Len() uint {
	return s.length
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}








func (s *SkipList) getPath(update []*node, key interface{}) *node {
	//
	// XXX
	//

	return nil
}













// XXX XXX XXX

/*
func New() {
	src := rand.NewSource(time.Now().UnixNano())

	src.Uint32()
}
*/

//func SetSeed()
