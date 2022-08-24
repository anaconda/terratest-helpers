# terratest-helpers

This module makes your terratest tests easier with pre-defined stages and other helpful functions.

All functions that init or apply the module allow you to specify an `errorFunc` which is called when the terraform command fails.

The error and the output of the terraform command are passed to that function so that you can define if the test should succeed or fail.

If you need to configure a terraform provider for your tests, save that configuration in `test/provider.tf`. The helper functions will automatically copy that configuration to the module directory before the test starts and delete it afterwards.

## Functions

* `StageSetupInitPlan`: Saves the `terraformOptions`, runs `terraform init` and `terraform plan`. You can specify an `errorFunc` here.
* `StageApply`: Runs `terraform apply`. You can specify an `errorFunc` here.
* `StageValidate`: If you want to check the state of the deployed infrastructure, you can do so in the validate stage. Pass a `validateFunc` into `StageValidate`.
* `StageDestroy`: Destroys the resources, cleans up test data and a provider configuration if it exists.

## Example

Will follow.
