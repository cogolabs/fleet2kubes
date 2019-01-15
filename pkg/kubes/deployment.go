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
				DNSPolicy string `json:"dnsPolicy" yaml:"dnsPolicy"`
				DNSConfig struct {
					Options []Option `json:"options"`
				} `json:"dnsConfig" yaml:"dnsConfig"`
				Containers []Container `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

func NewDeployment(name, image string, command []string, replicas, port int, env map[string]string,
	resources Resources) *Deployment {
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
	deploy.Spec.Template.Spec.DNSPolicy = "Default"
	deploy.Spec.Template.Spec.Containers = append(
		deploy.Spec.Template.Spec.Containers,
		Container{
			Name:      name,
			Image:     image,
			Command:   command,
			Ports:     []Port{{port}},
			Env:       newEnv(env),
			Resources: resources,
		},
	)
	return deploy
}
