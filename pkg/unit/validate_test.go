package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sj1 = `
[Service]
Restart=always
ExecStartPre=/bin/docker pull tech/sshjump:latest
ExecStartPre=-/bin/docker kill %p
ExecStartPre=-/bin/docker rm %p
ExecStart=/bin/docker run --rm --name=%p --net=none tech/sshjump:latest sshjump.git -proxyproto
ExecStartPost=/bin/pipework bond0 %p 172.16.0.11/19@172.16.31.250 @16
ExecStartPost=/bin/bash -c "nsenter -n -t $(docker inspect -f '{{.State.Pid}}' %p) -- ping -c 3 172.16.31.250"

[X-Fleet]
Conflicts=sshjump.*`
)

func Test1(t *testing.T) {
	u, err := NewUnit(sj1)
	assert.NoError(t, err)
	assert.NoError(t, u.Validate())
}
