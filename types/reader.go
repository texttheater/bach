package types

var ReaderType = readerType{}

type readerType struct {
}

func (t readerType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(readerType)
	return ok
}

func (t readerType) String() string {
	return "Reader"
}

func (t readerType) ElementType() Type {
	panic("Reader is not a sequence type")
}
