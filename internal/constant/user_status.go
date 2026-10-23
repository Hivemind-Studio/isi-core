package constant

type Status int

const (
	SUSPENDED Status = 0
	ACTIVE    Status = 1
	PENDING   Status = 2
	INACTIVE  Status = 3
)

var StatusMapping = map[string]Status{
	"ACTIVE":    ACTIVE,
	"SUSPENDED": SUSPENDED,
	"PENDING":   PENDING,
	"INACTIVE":  INACTIVE,
}

func GetStatusFromString(status string) Status {
	s := StatusMapping[status]
	return s
}
