package values

type Value interface {
	String() string
	Out() string
	Iter() <-chan Value
}
