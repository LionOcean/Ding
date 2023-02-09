package transfer

// some range slice s and judge whether filter could shot.
//
// s is any type(interface {}), filter func has three params and need one bool returned,
// current is value of slice piece, all is s self, index is index number of slice piece.
func some[T any](s []T, filter func(current T, all []T, index int) bool) bool {
	for i, v := range s {
		if filter(v, s, i) == true {
			return true
		}
	}
	return false
}

// includes range slice s and judge whether v is existed in s.
//
// Be sure that s piece is type comparable(string|int|unit|bool).
func includes[T comparable](s []T, v T) bool {
	return some(s, func(current T, all []T, index int) bool {
		return current == v
	})
}
