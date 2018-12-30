package kubes

type Service struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name   string `json:"name"`
		Labels struct {
			App string `json:"app"`
		} `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		LoadBalancerIP string        `json:"loadBalancerIP" yaml:"loadBalancerIP"`
		Ports          []ServicePort `json:"ports"`
		Selector       struct {
			App string `json:"app"`
		} `json:"selector"`
		SessionAffinity string `json:"sessionAffinity" yaml:"sessionAffinity"`
		Type            string `json:"type"`
	} `json:"spec"`
}

type ServicePort struct {
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort" yaml:"targetPort"`
}

func NewService(name string, IP string, port int) *Service {
	svc := &Service{
		APIVersion: "v1",
		Kind:       "Service",
	}
	svc.Metadata.Name = name
	svc.Metadata.Labels.App = name
	svc.Spec.LoadBalancerIP = IP
	svc.Spec.Ports = append(svc.Spec.Ports, ServicePort{"TCP", port, port})
	svc.Spec.Selector.App = name
	svc.Spec.SessionAffinity = "ClientIP"
	svc.Spec.Type = "LoadBalancer"
	return svc
}
