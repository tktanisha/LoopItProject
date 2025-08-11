package order_status

import "fmt"

type Status int

const (
	InUse Status = iota
	ReturnRequested
	Returned
)

var statusNames = [...]string{
	"In Use",
	"Return Requested",
	"Returned",
}

func (os Status) String() string {
	if os < 0 || int(os) >= len(statusNames) {
		return "unknown"
	}
	return statusNames[os]
}

func ParseStatus(statusStr string) (Status, error) {
	for i, v := range statusNames {
		if v == statusStr {
			return Status(i), nil
		}
	}
	return InUse, fmt.Errorf("invalid status: %s", statusStr)
}
