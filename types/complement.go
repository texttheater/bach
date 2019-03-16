package types

func Complement(a Type, b Type) Type {
	aUnion, ok := a.(unionType)
	if ok {
		var union unionType
		for _, disjunct := range aUnion {
			bUnion, ok := b.(unionType)
			if ok {
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
