package customconfig_test

import (
	"os"
	"testing"

	"github.com/ignite/apps/cs-client/customconfig"
	"github.com/stretchr/testify/assert"
)

func TestReadsCorrectly(t *testing.T) {
	content := `
version: 1
validation: sovereign
client:
  openapi:
    path: docs/static/openapi.yml
  cs-client:
    path: out/path
`

	file, _ := os.CreateTemp("", "*")
	defer file.Close()
	defer os.Remove(file.Name())
	file.WriteString(content)
	file.Close()

	c, err := customconfig.Read(file.Name())

	assert.NoError(t, err)
	assert.Equal(t, "out/path", c.Client.CsClient.Path)
}
