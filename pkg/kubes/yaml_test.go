package kubes

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var expected = `apiVersion: v1
kind: Service
metadata:
  name: test1
  labels:
    app: test1
spec:
  loadBalancerIP: 1.2.3.4
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
  selector:
    app: test1
  sessionAffinity: ClientIP
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test1
  labels:
    app: test1
spec:
  replicas: 2
  selector:
    matchLabels:
      app: test1
  template:
    metadata:
      labels:
        app: test1
    spec:
      dnsConfig:
        options:
        - name: ndots
          value: "1"
      containers:
      - name: test1
        image: httpd
        command:
        - httpd
        - -listen
        - :80
        ports:
        - containerPort: 80
        env:
        - name: FOO
          value: BAR
`

func TestYAML(t *testing.T) {
	output := bytes.NewBufferString("")
	svc := NewService("test1", "1.2.3.4", 80)
	dpl := NewDeployment("test1", "httpd", []string{"httpd", "-listen", ":80"}, 2, 80, map[string]string{"FOO": "BAR"})

	err := yaml.NewEncoder(output).Encode(svc)
	assert.NoError(t, err)
	fmt.Fprint(output, "---\n")
	err = yaml.NewEncoder(output).Encode(dpl)
	assert.NoError(t, err)

	assert.Equal(t, expected, output.String())
}
