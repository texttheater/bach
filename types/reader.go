package types

type ReaderType struct {
}

func (t ReaderType) Subsumes(other Type) bool {
	if (VoidType{}).Subsumes(other) {
		return true
	}
	_, ok := other.(ReaderType)
	return ok
}

func (t ReaderType) String() string {
	return "Reader"
}

func (t ReaderType) ElementType() Type {
	panic("Reader is not a sequence type")
}
