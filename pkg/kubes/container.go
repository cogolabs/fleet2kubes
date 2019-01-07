package kubes

type Container struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command"`
	Ports   []Port   `json:"ports" yaml:"ports,omitempty"`
	Env     Env      `json:"env" yaml:"env,omitempty"`
}

type Port struct {
	ContainerPort int `json:"containerPort" yaml:"containerPort,omitempty"`
}

type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Env []Option

func newEnv(envMap map[string]string) Env {
	var env Env
	for name, val := range envMap {
		env = append(env, Option{Name: name, Value: val})
	}

	return env
}
