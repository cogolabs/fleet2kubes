package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitBuilder(t *testing.T) {
	raw := `[Service]
Restart=always
RestartSec=60s
ExecStartPre=/bin/docker pull tools/builder:v3.1
ExecStartPre=-/bin/docker kill %p%i
ExecStartPre=-/bin/docker rm %p%i
ExecStart=/bin/docker run --rm --name %p%i --net=host -e TZ=America/New_York -e SERVER=wss://docker.cogolo.net -v /var/run/docker.sock:/var/run/docker.sock -v /var/run/nscd:/var/run/nscd -v /etc/pki:/etc/pki:ro tools/builder:v3.1
ExecStartPost=/bin/pipework bond0 %p 172.16.1.126/19@172.16.0.1 @16
ExecStop=-/bin/docker kill %p%i`
	u, err := NewUnit(raw)
	assert.NoError(t, err)
	assert.Equal(t, "tools/builder:v3.1", u.RunImage)
	assert.Equal(t, []string{}, u.RunCommand)
	assert.Equal(t, map[string]string{"bond0.16": "172.16.1.126/19"}, u.Network)

}
