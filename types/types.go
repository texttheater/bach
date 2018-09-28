package types

type Type interface {
	Subsumes(Type) bool
}

///////////////////////////////////////////////////////////////////////////////

type AnyType struct {
}

func (self AnyType) Subsumes(other Type) bool {
	return true
}

///////////////////////////////////////////////////////////////////////////////

type SeqType struct {
	ElementType Type
}

func (self SeqType) Subsumes(other Type) bool {
	otherSeqType, ok := other.(SeqType)
	if !ok {
		return false
	}
	return self.ElementType.Subsumes(otherSeqType.ElementType)
}

///////////////////////////////////////////////////////////////////////////////

type NullType struct {
}

func (self NullType) Subsumes(other Type) bool {
	_, ok := other.(NullType)
	return ok
}

///////////////////////////////////////////////////////////////////////////////

type BooleanType struct {
}

func (self BooleanType) Subsumes(other Type) bool {
	_, ok := other.(NullType)
	return ok
}

///////////////////////////////////////////////////////////////////////////////

type NumberType struct {
}

func (self NumberType) Subsumes(other Type) bool {
	_, ok := other.(NumberType)
	return ok
}

///////////////////////////////////////////////////////////////////////////////

type StringType struct {
}

func (self StringType) Subsumes(other Type) bool {
	_, ok := other.(StringType)
	return ok
}

///////////////////////////////////////////////////////////////////////////////

type ArrayType struct {
	ElementType Type
	HasLength   bool
	Length      uint
}

func (self ArrayType) Subsumes(other Type) bool {
	otherArrayType, ok := other.(ArrayType)
	if !ok {
		return false
	}
	if self.HasLength {
		if !otherArrayType.HasLength {
			return false
		}
		if otherArrayType.Length != self.Length {
			return false
		}
	}
	return self.ElementType.Subsumes(otherArrayType.ElementType)
}

///////////////////////////////////////////////////////////////////////////////

type ObjectType struct {
	fieldTypeMap map[string]Type
}

func (self ObjectType) Subsumes(other Type) bool {
	otherObjectType, ok := other.(ObjectType)
	if !ok {
		return false
	}
	for field, b := range otherObjectType.fieldTypeMap {
		if a, ok := self.fieldTypeMap[field]; ok {
			if !a.Subsumes(b) {
				return false
			}
		}
	}
	return true
}
