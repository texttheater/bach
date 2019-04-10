package types

func Complement(a Type, b Type) Type {
	if aUnion, ok := a.(unionType); ok {
		var union unionType
		for _, disjunct := range aUnion {
			if bUnion, ok := b.(unionType); ok {
				for _, bDisjunct := range bUnion {
					disjunct = Complement(disjunct, bDisjunct)
				}
			} else {
				disjunct = Complement(disjunct, b)
			}
			union = typeAppend(union, disjunct)
		}
	}
	if b.Subsumes(a) {
		return VoidType
	}
	return a
}
