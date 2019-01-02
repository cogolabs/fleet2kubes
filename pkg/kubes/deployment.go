package kubes

type Deployment struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name   string `json:"name"`
		Labels struct {
			App string `json:"app"`
		} `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
		Selector struct {
			MatchLabels struct {
				App string `json:"app"`
			} `json:"matchLabels" yaml:"matchLabels"`
		} `json:"selector"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `json:"app"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				DNSConfig struct {
					Options []Option `json:"options"`
				} `json:"dnsConfig" yaml:"dnsConfig"`
				Containers []Container `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

type Container struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command"`
	Ports   []Port   `json:"ports"`
}

type Port struct {
	ContainerPort int `json:"containerPort" yaml:"containerPort"`
}

type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewDeployment(name, image string, command []string, replicas, port int) *Deployment {
	deploy := &Deployment{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	}
	deploy.Metadata.Name = name
	deploy.Metadata.Labels.App = name
	deploy.Spec.Replicas = replicas
	deploy.Spec.Selector.MatchLabels.App = name
	deploy.Spec.Template.Metadata.Labels.App = name
	deploy.Spec.Template.Spec.DNSConfig.Options = append(
		deploy.Spec.Template.Spec.DNSConfig.Options,
		Option{"ndots", "1"},
	)
	deploy.Spec.Template.Spec.Containers = append(
		deploy.Spec.Template.Spec.Containers,
		Container{
			Name:    name,
			Image:   image,
			Command: command,
			Ports:   []Port{{port}},
		},
	)
	return deploy
}
