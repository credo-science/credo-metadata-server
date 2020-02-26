package event

type Type int

type Metadata map[string]string

const (
	Detection Type = iota
	Ping
)

var stringToId = map[string]Type{
	"detection": Detection,
	"ping":      Ping,
}

func (t Type) String() string {
	return [...]string{"detection", "ping"}[t]
}

func StringToId(t string) (Type, bool) {
	val, ok := stringToId[t]
	return val, ok
}
