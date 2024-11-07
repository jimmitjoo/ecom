package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwaggerValidation(t *testing.T) {
	err := ValidateSwaggerSpec("../../../docs/swagger.yaml")
	assert.NoError(t, err, "Swagger spec should be valid")
}
