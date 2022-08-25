# terratest-helpers

This module makes your terratest tests easier with pre-defined stages and other helpful functions.

All functions that init or apply the module allow you to specify an `errorFunc` which is called when the terraform command fails.

The error and the output of the terraform command are passed to that function so that you can define if the test should succeed or fail.

If you need to configure a terraform provider for your tests, save that configuration in `test/provider.tf`. The helper functions will automatically copy that configuration to the module directory before the test starts and delete it afterwards.

## Functions

### Stages

When stages are used, you can set a `SKIP_$stage` environment variable to skip a certain stage, e.g. `SKIP_destroy`. This makes developing and debugging tests much easier.

* `StageSetup`: Saves the `terraformOptions`, runs `terraform init` and `terraform plan`. You can specify an `errorFunc` here.
* `StageApply`: Runs `terraform apply`. You can specify an `errorFunc` here.
* `StageValidate`: If you want to check the state of the deployed infrastructure, you can do so in the validate stage. Pass a `validateFunc` into `StageValidate`.
* `StageDestroy`: Runs `terraform destroy` and calls the `Cleanup` function.

### Others

* `Cleanup`: Removes the test data directory `.test-data` and test provider configuration `test-provider.tf`.

## Example

### Standard test

```go
func TestS3BucketDefault(t *testing.T) {
	defer helpers.StageDestroy(t, TerraformDir)

	// set everything up for the terraform apply later
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: TerraformDir,
		Vars: map[string]interface{}{
			"name":        uuid.New(),
		},
	})

	helpers.StageSetup(t, TerraformDir, terraformOptions)
	helpers.StageApply(t, TerraformDir)
}
```

### Expected to fail

For a test that is expected to fail on plan, use the `Cleanup` function.

```go
func TestS3BucketBackupOnVersioningOff(t *testing.T) {
	defer helpers.Cleanup(t, TerraformDir)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: TerraformDir,
		// NoColor is important for the string.Contains later
		NoColor: true,
		Vars: map[string]interface{}{
			"name":        uuid.New(),
			"backup":      "on",
			"versioning":  "Disabled",
		},
	})

	// The terraform plan is expected to fail as versioning can't be turned off with backups turned on.
	helpers.StageSetup(t, TerraformDir, terraformOptions, func(err error, stdoutStderr string) {
		// This error is expected, the precondition failed
		if err != nil {
			if !strings.Contains(stdoutStderr, "Error: Resource precondition failed") || !strings.Contains(stdoutStderr, "Versioning cannot be disabled when backups are enabled") {
				assert.Fail(t, "The precondition checking versioning and backup configuration did not fail, but there is another error.")
			}
		} else {
			// There must be an error in the precondition
			assert.Fail(t, "There are no errors, but the precondition checking versioning and backup configuration must fail.")
		}
	})
}
```

## Development

When cloning this repository, run `make init` locally. This sets up pre-commit, which takes care of the go formatting.

## Versioning

Versioning follows the go standard: Semantic versioning with a `v` prefix. Tags are currently created and pushed manually.
