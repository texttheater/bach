package types

type ReaderType struct {
}

func (t ReaderType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case ReaderType:
		return true
	default:
		return false
	}
}

func (t ReaderType) String() string {
	return "Reader"
}

func (t ReaderType) ElementType() Type {
	panic("Reader is not a sequence type")
}
