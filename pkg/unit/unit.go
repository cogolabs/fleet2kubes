package unit

import (
	unit "github.com/cogolabs/fleet2kubes/pkg/fleet-unit"
)

type Unit struct {
	*unit.UnitFile

	RunFlags   map[string]string
	RunImage   string
	RunCommand []string

	Network map[string]string

	Macros map[string]string
}

func NewUnit(raw string) (*Unit, error) {
	uf, err := unit.NewUnitFile(raw)
	if err != nil {
		return nil, err
	}
	u := &Unit{UnitFile: uf}
	u.Network = u.parseNetwork()
	u.RunImage, u.RunCommand, u.RunFlags = u.parseExecStart()
	return u, nil
}

func (u *Unit) FirstValue(section, name string) string {
	if u == nil || u.Options == nil {
		return ""
	}
	values := u.Contents[section][name]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
