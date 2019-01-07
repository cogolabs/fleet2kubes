package kubes

type Container struct {
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Command   []string  `json:"command"`
	Ports     []Port    `json:"ports" yaml:"ports,omitempty"`
	Env       Env       `json:"env" yaml:"env,omitempty"`
	Resources Resources `json:"resources" yaml:"resources,omitempty"`
}

type Port struct {
	ContainerPort int `json:"containerPort" yaml:"containerPort,omitempty"`
}

type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Env []Option

type Resources struct {
	Requests struct {
		Memory string `json:"memory"`
		CPU    string `json:"cpu"`
	} `json:"requests"`
	Limits struct {
		Memory string `json:"memory"`
		CPU    string `json:"cpu"`
	} `json:"limits"`
}

func newEnv(envMap map[string]string) Env {
	var env Env
	for name, val := range envMap {
		env = append(env, Option{Name: name, Value: val})
	}

	return env
}
