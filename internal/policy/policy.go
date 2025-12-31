package policy

type ProtectionLevel int

const (
	Allowed ProtectionLevel = iota
	ConfirmRequired
	Denied
)

func (p ProtectionLevel) String() string {
	switch p {
	case Allowed:
		return "ALLOWED"
	case ConfirmRequired:
		return "CONFIRM_REQUIRED"
	case Denied:
		return "DENIED"
	default:
		return "UNKNOWN"
	}
}

type Decision struct {
	Level  ProtectionLevel
	Reason string
}

