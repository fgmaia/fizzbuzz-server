package apps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockApps(t *testing.T) {
	customApp := &FizzbuzzApp{}
	MockTestsApp(customApp)
	assert.NotNil(t, App)
}
