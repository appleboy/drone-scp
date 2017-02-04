package easyssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyFile(t *testing.T) {
	// missing file
	_, err := getKeyFile("abc")
	assert.Error(t, err)
	assert.Equal(t, "open abc: no such file or directory", err.Error())

	// wrong format
	_, err = getKeyFile("../tests/.ssh/id_rsa.pub")
	assert.Error(t, err)
	assert.Equal(t, "ssh: no key found", err.Error())

	_, err = getKeyFile("../tests/.ssh/id_rsa")
	assert.NoError(t, err)
}
