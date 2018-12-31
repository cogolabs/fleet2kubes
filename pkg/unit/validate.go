package unit

import (
	"fmt"
	"strings"
)

func (u *Unit) ValidateAll() error {
	if err := u.ValidateKeys(); err != nil {
		return err
	}
	if err := u.ValidateExecStartPre(); err != nil {
		return err
	}
	if err := u.ValidateExecStart(); err != nil {
		return err
	}
	if err := u.ValidateExecStartPost(); err != nil {
		return err
	}
	return nil
}

func (u *Unit) ValidateKeys() error {
	serviceKeys := u.Contents["Service"]
	for k := range serviceKeys {
		switch k {
		case
			"Restart",
			"ExecStartPre",
			"ExecStart",
			"ExecStartPost":
			continue
		default:
			return fmt.Errorf("unexpected service key: %+v", k)
		}
	}
	return nil
}

func (u *Unit) ValidateExecStartPre() error {
	for _, elt := range u.Contents["Service"]["ExecStartPre"] {
		s := stripExec(elt)
		switch {
		case
			strings.HasPrefix(s, "docker pull"),
			strings.HasPrefix(s, "docker kill"),
			strings.HasPrefix(s, "docker rm"):
			continue
		default:
			return fmt.Errorf("unexpected ExecStartPre: %+v", elt)
		}
	}
	return nil
}

func (u *Unit) ValidateExecStart() error {
	for _, elt := range u.Contents["Service"]["ExecStart"] {
		s := stripExec(elt)
		switch {
		case
			strings.HasPrefix(s, "docker run"):
			continue
		default:
			return fmt.Errorf("unexpected ExecStart: %+v", elt)
		}
	}
	return nil
}

func (u *Unit) ValidateExecStartPost() error {
	for _, elt := range u.Contents["Service"]["ExecStartPost"] {
		s := stripExec(elt)
		switch {
		case
			strings.HasPrefix(s, "pipework"),
			strings.HasPrefix(s, "bash -c \"nsenter"):
			continue
		default:
			return fmt.Errorf("unexpected ExecStartPost: %+v", elt)
		}
	}
	return nil
}
