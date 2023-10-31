package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func DefaultOptions(t *testing.T, terraformOptions *terraform.Options) *terraform.Options {
	// Check if a variable file with the name of the current test exists and use it
	currentTestVarFile := filepath.Join("variables", fmt.Sprintf("%s.tfvars", t.Name()))
	if files.FileExists(currentTestVarFile) {
		terraformOptions.VarFiles = append(terraformOptions.VarFiles, filepath.Join("test", currentTestVarFile))
	}

	// By default, our terraform module is in the git repository root and our tests in the 'test' directory
	// Therefore, the we set the TerraformDir to ".." if none is specified explicitly
	if terraformOptions.TerraformDir == "" {
		path, err := os.Getwd()
		if err != nil {
			assert.FailNow(t, "Could not get current working directory, aborting test")
		}
		terraformOptions.TerraformDir = filepath.Join(path, "../")
	}

	// This enables parallel testing by ensuring that every TerraformDir lives in its own temporary path
	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, terraformOptions.TerraformDir, ".")
	terraformOptions.TerraformDir = tempTestFolder

	return terraform.WithDefaultRetryableErrors(t, terraformOptions)
}
