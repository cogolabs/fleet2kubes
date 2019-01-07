package kubes

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var (
	expected = `apiVersion: v1
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
	expectedCronJob = `apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: test3
  annotations:
    description: A test cron job
    documentation: http://git.colofoo.net/fleet/test3
spec:
  concurrencyPolicy: Forbid
  schedule: '* 7 * 1 *'
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
          - name: test3
            image: cleanup
            command:
            - /bin/cleanup
            - -f
            env:
            - name: FOO
              value: BAR
            resources:
              requests:
                memory: 4Gi
                cpu: "1.5"
              limits:
                memory: 16Gi
                cpu: "4"
`
)

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

func TestCronJob(t *testing.T) {
	output := bytes.NewBufferString("")
	env := map[string]string{"FOO": "BAR"}
	resources := Resources{}
	resources.Limits.Memory = "16Gi"
	resources.Limits.CPU = "4"
	resources.Requests.Memory = "4Gi"
	resources.Requests.CPU = "1.5"
	annotations := Annotations{"description": "A test cron job", "documentation": "http://git.colofoo.net/fleet/test3"}
	cj := NewCronJob("test3", "*-01-* 07:*", "Forbid", "OnFailure", "cleanup", []string{"/bin/cleanup", "-f"}, env, resources, annotations)
	err := yaml.NewEncoder(output).Encode(cj)
	assert.NoError(t, err)
	assert.Equal(t, expectedCronJob, output.String())
}

func TestParseSchedule(t *testing.T) {
	tests := []struct {
		schedule string
		expected string
	}{
		{"14:35", "35 14 * * *"},
		{"Fri 6:*", "* 6 * * 5"},
		{"*-*-* *:22", "22 * * * *"},
		{"Tue 2019-01-15 2:44", "44 2 15 1 2"},
	}

	for _, tt := range tests {
		schedule := parseSchedule(tt.schedule)
		assert.Equal(t, tt.expected, schedule, "For: %s", tt.schedule)
	}
}
