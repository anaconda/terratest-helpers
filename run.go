package helpers

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func RunValidate(t *testing.T, terraformOptions terraform.Options, validateFunc func()) {
	RunOptionsValidate(t, &terraform.Options{}, validateFunc)
}

func RunNoValidate(t *testing.T) {
	RunOptionsValidate(t, &terraform.Options{}, func() {})
}

// RunOptionsNoValidate runs applies and destroys the module with the configured
// variables.
// Use this function when for tests where you expect terraform to run successfully
// without any additional validation.
func RunOptionsNoValidate(t *testing.T, terraformOptions *terraform.Options) {
	RunOptionsValidate(t, terraformOptions, func() {})
}

func RunOptionsValidate(t *testing.T, terraformOptions *terraform.Options, validateFunc func()) {
	tOptions := *DefaultOptions(t, terraformOptions)

	defer StageDestroy(t, tOptions.TerraformDir)
	StageSetup(t, tOptions.TerraformDir, &tOptions)
	StageApply(t, tOptions.TerraformDir)
	StageValidate(t, validateFunc)
}
