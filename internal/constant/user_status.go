package constant

type Status int

const (
	ACTIVE    Status = 1
	SUSPENDED Status = 0
)

var StatusMapping = map[string]Status{
	"ACTIVE":    ACTIVE,
	"SUSPENDED": SUSPENDED,
}

func GetStatusFromString(status string) Status {
	s := StatusMapping[status]
	return s
}
