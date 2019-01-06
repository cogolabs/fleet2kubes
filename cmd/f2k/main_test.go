package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployService(t *testing.T) {
	actual := bytes.NewBufferString("")
	err := do("test2.service", actual)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile("test2.yaml")
	assert.NoError(t, err)

	assert.Equal(t, string(expected), actual.String())
}

func TestCronJob(t *testing.T) {
	actual := bytes.NewBufferString("")
	err := do("test4.service", actual)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile("test4.yaml")
	assert.NoError(t, err)

	assert.Equal(t, string(expected), actual.String())
}
