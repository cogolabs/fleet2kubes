package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	test1 = `
[Service]
Restart=always
ExecStartPre=/bin/docker pull tech/test1:latest
ExecStartPre=-/bin/docker kill %p
ExecStartPre=-/bin/docker rm %p
ExecStart=/bin/docker run --rm --name=%p --net=none tech/test1:latest test1.git -test1
ExecStartPost=/bin/pipework bond0 %p 172.16.0.11/19@172.16.31.250 @16
ExecStartPost=/bin/bash -c "nsenter -n -t $(docker inspect -f '{{.State.Pid}}' %p) -- ping -c 3 172.16.31.250"

[X-Fleet]
Conflicts=test2.*`
)

func Test1(t *testing.T) {
	u, err := NewUnit(test1)
	assert.NoError(t, err)
	assert.NoError(t, u.Validate())
}
