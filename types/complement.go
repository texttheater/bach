package types

func Complement(a Type, b Type) Type {
	if aUnion, ok := a.(UnionType); ok {
		var union UnionType
		for _, disjunct := range aUnion {
			if bUnion, ok := b.(UnionType); ok {
				for _, bDisjunct := range bUnion {
					disjunct = Complement(disjunct, bDisjunct)
				}
			} else {
				disjunct = Complement(disjunct, b)
			}
			union = typeAppend(union, disjunct)
		}
		if len(union) == 1 {
			return union[0]
		}
		return union
	}
	if b.Subsumes(a) {
		return VoidType{}
	}
	return a
}
