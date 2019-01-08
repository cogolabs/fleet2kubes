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
		Memory string `json:"memory" yaml:"memory,omitempty"`
		CPU    string `json:"cpu" yaml:"cpu,omitempty"`
	} `json:"requests" yaml:"requests,omitempty"`
	Limits struct {
		Memory string `json:"memory" yaml:"memory,omitempty"`
		CPU    string `json:"cpu" yaml:"cpu,omitempty"`
	} `json:"limits" yaml:"limits,omitempty"`
}

func newEnv(envMap map[string]string) Env {
	var env Env
	for name, val := range envMap {
		env = append(env, Option{Name: name, Value: val})
	}

	return env
}

func NewResources(memLimit, cpuLimit, reqMemLimit, reqCPULimit string) Resources {
	resources := Resources{}
	if memLimit != "" {
		resources.Limits.Memory = memLimit
	}
	if cpuLimit != "" {
		resources.Limits.CPU = cpuLimit
	}
	if reqMemLimit != "" {
		resources.Requests.Memory = reqMemLimit
	}
	if reqCPULimit != "" {
		resources.Requests.CPU = reqCPULimit
	}

	return resources
}
