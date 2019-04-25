package types

func TupType(elementTypes ...Type) Type {
	t := VoidArrType
	for i := len(elementTypes) - 1; i >= 0; i-- {
		t = &NearrType{elementTypes[i], t}
	}
	return t
}
