package buyer_request_status

import "fmt"

type Status int

const (
	Pending Status = iota
	Approved
	Rejected
)

var statusNames = [...]string{
	"Pending",
	"Approved",
	"Rejected",
}

func (s Status) String() string {
	if s < 0 || int(s) >= len(statusNames) {
		return "unknown"
	}
	return statusNames[s]
}

func ParseStatus(statusStr string) (Status, error) {
	for i, v := range statusNames {
		if v == statusStr {
			return Status(i), nil
		}
	}
	return Pending, fmt.Errorf("invalid status: %s", statusStr)
}
