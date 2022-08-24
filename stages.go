package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// StageSetupAndInit is a helper function that sets up the environment to run validation
// It performs the following steps:
// 0. Copies provider configuration for the AWS provider
// 1. Saving the terraformOptions for later stages
// 2. Runs `terraform init` and `terraform plan`
// 3. Verifies that the plan is not running in the remote backend.
// 4. Executes a user specified function to validate any errors that might have occurred.
func StageSetupInitPlan(t *testing.T, terraformDir string, terraformOptions *terraform.Options, errorFunc ...func(err error, stdoutStderr string)) {
	// We only allow 1 errorFunc
	if len(errorFunc) > 1 {
		assert.FailNow(t, "You must specify exactly zero or one errorFunc's")
	}

	// This copies the provider configuration file if it exists.
	// Some modules don't need a provider.tf for the tests, therefore we check if it exists
	if files.FileExists("provider.tf") {
		err := files.CopyFile("provider.tf", filepath.Join(terraformDir, "test-provider.tf"))
		if err != nil {
			assert.FailNow(t, "Could not copy provider.tf file, aborting test")
		}
	}

	test_structure.RunTestStage(t, "setup", func() {
		test_structure.SaveTerraformOptions(t, terraformDir, terraformOptions)

		stdoutStderr, err := terraform.InitAndPlanE(t, terraformOptions)

		// Verify we're running locally to not interfere with non-testing infrastructure
		exit := assert.NotContains(t, stdoutStderr, "backend \"remote\"", "Running plan in the remote backend", "Plan is running in remote backend")

		if !exit {
			assert.FailNow(t, "terraform is configured with remote backend")
		}

		// Run the errorFunc if specified, else fail
		if err != nil && errorFunc != nil {
			errorFunc[0](err, stdoutStderr)
		} else if err != nil {
			assert.FailNow(t, err.Error())
		}
	})
}

// StageApply runs "terraform apply" and validates any errors against the user specified errorFunc.
func StageApply(t *testing.T, terraformDir string, errorFunc ...func(err error, stdoutStderr string)) {
	// We only allow 1 errorFunc
	if len(errorFunc) > 1 {
		assert.FailNow(t, "You must specify exactly zero or one errorFunc's")
	}

	test_structure.RunTestStage(t, "apply", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, terraformDir)
		stdoutStderr, err := terraform.ApplyE(t, terraformOptions)

		// Run the errorFunc if specified, else fail
		if err != nil && errorFunc != nil {
			errorFunc[0](err, stdoutStderr)
		} else if err != nil {
			assert.FailNow(t, err.Error())
		}
	})
}

// StageValidate runs the user specified validation function to ensure that the deployed infrastructure
// is configured as it should be.
func StageValidate(t *testing.T, validateFunc ...func()) {
	// We only allow 1 validateFunc
	if len(validateFunc) > 1 {
		assert.FailNow(t, "You must specify exactly zero or one validateFunc's")
	}

	// If a validateFunc is specified, run it. If not, do nothing
	if validateFunc != nil {
		test_structure.RunTestStage(t, "validate", validateFunc[0])
	}
}

// StageDestroy runs "terraform destroy", cleans up the test data and removes the provider configuration.
func StageDestroy(t *testing.T, terraformDir string) {
	test_structure.RunTestStage(t, "destroy", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, terraformDir)
		terraform.Destroy(t, terraformOptions)
		test_structure.CleanupTestDataFolder(t, terraformDir)

		// Clean up the test-provider.tf
		providerPath := filepath.Join(terraformDir, "test-provider.tf")
		if files.FileExists(providerPath) {
			os.Remove(providerPath)
		}
	})
}
