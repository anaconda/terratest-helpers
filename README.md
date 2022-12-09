# terratest-helpers

This module makes your terratest tests easier with pre-defined stages and other helpful functions.

All functions that init or apply the module allow you to specify an `errorFunc` which is called when the terraform command fails.

The error and the output of the terraform command are passed to that function so that you can define if the test should succeed or fail.

## Repository layout

To ensure correct function of this module, you need to:

- Put your tests in the directory `test`.
- Put variable files for your tests in `test/variables` and name them as the test is named
- Put provider configuration into the `test/provider.tf` file

## Functions

### Run functions

Those functions are the easy start into testing, with more complex scenarios covered below. They are called “run functions” because their name starts with `Run`.

All run functions destroy the deployed infrastructure afterwards.
Run functions abstract [stages](#stages) for you.

To only run some stages, set a `SKIP_$stage` environment variable to skip a certain stage, e.g. `SKIP_destroy`. This makes developing and debugging tests much easier.

For examples, see [the examples section](#examples).

Available run functions are:

- `RunNoValidate`: Run with default options and no extra validation
- `RunValidate`: Same as `RunNoValidate`, but with a validation function
- `RunOptionsNoValidate`: Run with configured terraform options, no extra validation
- `RunOptionsValidate`: Same as `RunOptionsNoValidate`, but with a validation function

In a validation function, you can inspect the deployed infrastructure and have the test fail with `assert.Fail` if it does not match what you are expecting.

### Terraform Options

You can use the `DefaultOptions` function for default terraformOptions to be set automatically. Those are:

- `TerraformDir`, set to the absolute path of the module as the tests are in a subdirectory named `tests`
- `VarFiles` is appended the with the `.tfvars` file in `test/variables` that has the same as the running test

### Stages

:information_source: You need stages if you expect e.g. the `terraform plan` to fail for a test and want to validate the errors that occur, see the examples below.

Available stages are:

- `setup`: Saves the `terraformOptions`, runs `terraform init` and `terraform plan`. You can specify an `errorFunc` here.
- `apply`: Runs `terraform apply`. You can specify an `errorFunc` here.
- `validate`: If you want to check the state of the deployed infrastructure, you can do so in the validate stage. Pass a `validateFunc` into `StageValidate`.
- `destroy`: Runs `terraform destroy` and calls the `Cleanup` function.

### Others

- `Cleanup`: Removes the test data directory `.test-data` and test provider configuration `test-provider.tf`.

## Examples

### Standard test

```go
func TestDefaults(t *testing.T) {
	helpers.RunNoValidate(t)
}
```

with this, the test will run the module with the variables in `test/variables/TestDefaults.tfvars` and expect all operations to be successful.

### Specifying dynamic variables

When you need to define not only constant variables, but also dynamic ones, you can use `RunOptionsNoValidate` as follows:

```go
func TestWithVars(t *testing.T) {
	helpers.RunOptionsNoValidate(t, &terraform.Options{
		Vars: map[string]interface{}{
			"name": uuid.New(),
		},
	})
}
```

This is needed for e.g. tests with S3 buckets to ensure every bucket is named uniquely.

### Expected to fail on apply

If a test is expected to fail on apply, defer the `destroy` stage and pass an error function to the `apply` stage:

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
	helpers.StageApply(t, TerraformDir, func(err error, stdoutStderr string) {
		// This error is expected
		if err != nil {
			if !strings.Contains(stdoutStderr, "Error: Some expected error") {
				assert.Fail(t, "We expected an error, but it did not occur. Instead, another error was returned.")
			}
		}
	})
}
```

### Expected to fail on plan

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
