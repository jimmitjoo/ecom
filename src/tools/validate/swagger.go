package validate

import (
	"fmt"
	"os/exec"
)

func ValidateSwaggerSpec(specPath string) error {
	cmd := exec.Command("swagger", "validate", specPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("swagger validation failed: %s\n%s", err, output)
	}
	return nil
}
