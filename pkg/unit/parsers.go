package unit

import (
	"strings"

	"github.com/labstack/gommon/log"
)

func (u *Unit) parseExecStart() (string, []string, map[string]string) {
	v := u.FirstValue("Service", "ExecStart")
	f := &FlagSet{}
	g := strings.Split(v, "docker run ")
	if len(g) <= 1 {
		return "", nil, nil
	}
	argv := g[1]
	for strings.Contains(argv, "\t") {
		argv = strings.Replace(argv, "\t", " ", -1)
	}
	for strings.Contains(argv, "  ") {
		argv = strings.Replace(argv, "  ", " ", -1)
	}
	err := f.Parse(strings.Split(argv, " "))
	if err != nil {
		log.Warn(err)
		return "", nil, nil
	}
	args := f.Args()
	if len(args) < 1 {
		return "", nil, nil
	}
	s := args[0]
	if !strings.Contains(s, ":") {
		s += ":latest"
	}
	return s, args[1:], f.Values()
}

func (u *Unit) parseNetwork() map[string]string {
	r := map[string]string{}
	for _, opt := range u.Contents["Service"]["ExecStartPost"] {
		g := strings.Split(opt, "bin/pipework ")
		if len(g) < 2 {
			continue
		}
		v := strings.Split(g[1], " ")
		if len(v) < 3 {
			continue
		}
		vlan := ""
		if len(v) > 3 && strings.Contains(v[len(v)-1], "@") {
			vlan = "." + strings.Split(v[len(v)-1], "@")[1]
		}
		if v[1] == "-i" && len(v) > 4 {
			r[v[0]+vlan] = stripNetwork(u.applyMacros(v[4]))
			continue
		}
		r[v[0]+vlan] = stripNetwork(u.applyMacros(v[2]))
	}
	return r
}

func (u *Unit) applyMacros(s string) string {
	for k, v := range u.Macros {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

func stripExec(s string) string {
	if s == "" {
		return s
	}
	if s[0] == '-' {
		s = s[1:]
	}
	if s[0] == '/' {
		s = s[1:]
	}
	if strings.HasPrefix(s, "usr/") {
		s = s[4:]
	}
	if strings.HasPrefix(s, "bin/") {
		s = s[4:]
	}
	return s
}

func stripNetwork(s string) string {
	return strings.Split(s, "@")[0]
}
