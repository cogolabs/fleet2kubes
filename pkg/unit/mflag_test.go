package unit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunFlags(t *testing.T) {
	f := FlagSet{}
	f.Parse([]string{"-net=none", "-e", "FOO=BAR", "centos:latest", "bash"})

	assert.Equal(t, map[string]string{"net": "none"}, f.Values())
	assert.Equal(t, map[string]string{"FOO": "BAR"}, f.Env())
	assert.Equal(t, []string{"centos:latest", "bash"}, f.Args())
}

var test2 = `--rm --name %p-%i --net=none -e SERVER=wss://registry.colofoo.net -v /var/run/docker.sock:/var/run/docker.sock -v /var/run/nscd:/var/run/nscd -v /etc/pki:/etc/pki:ro -v /bin:/bin:ro -v /lib64:/lib64:ro docker.colofoo.net/coreos-registry-build-worker:latest /bin/bash -c '/bin/pipework --wait; /quay-builder'`

func TestRunFlags2(t *testing.T) {
	f := FlagSet{}
	f.Parse(strings.Split(test2, " "))

	assert.Equal(t, map[string]string{"-rm": "", "-name": "%p-%i", "-net": "none", "v": "/lib64:/lib64:ro"}, f.Values())
	assert.Equal(t, map[string]string{"SERVER": "wss://registry.colofoo.net"}, f.Env())
	assert.Equal(t, []string{"docker.colofoo.net/coreos-registry-build-worker:latest", "/bin/bash", "-c", "'/bin/pipework", "--wait;", "/quay-builder'"}, f.Args())
}
