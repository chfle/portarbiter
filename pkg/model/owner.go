package model

type OwnerType string

const (
	OwnerProcess OwnerType = "process"
	OwnerSystemd OwnerType = "systemd"
	OwnerDocker  OwnerType = "docker"
	OwnerCompose OwnerType = "docker-compose"
)

// String implements fmt.Stringer
func (o OwnerType) String() string {
	return string(o)
}

type Owner interface {
	Type() OwnerType
	ID() string
	Describe() string
	Kill(force bool) error
}

