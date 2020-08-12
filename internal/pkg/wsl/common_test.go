package wsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWslParsing(t *testing.T) {
	input := `  NAME                   STATE           VERSION
	* Ubuntusl               Running         2
	  Ubuntu-18.04           Running         1
	  Ubuntu-20.04           Stopped         2
	`

	distros, err := parseDistroOutput(input)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(distros))

	if len(distros) >= 1 {
		distro := distros[0]
		assert.Equal(t, true, distro.IsDefault)
		assert.Equal(t, "Ubuntusl", distro.Name)
		assert.Equal(t, "Running", distro.State)
		assert.Equal(t, "2", distro.Version)
	}

	if len(distros) >= 2 {
		distro := distros[1]
		assert.Equal(t, false, distro.IsDefault)
		assert.Equal(t, "Ubuntu-18.04", distro.Name)
		assert.Equal(t, "Running", distro.State)
		assert.Equal(t, "1", distro.Version)
	}

	if len(distros) >= 3 {
		distro := distros[2]
		assert.Equal(t, false, distro.IsDefault)
		assert.Equal(t, "Ubuntu-20.04", distro.Name)
		assert.Equal(t, "Stopped", distro.State)
		assert.Equal(t, "2", distro.Version)
	}

}
