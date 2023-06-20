//go:build integration
// +build integration

package test

import (
	"github.com/gruntwork-io/terratest/modules/random"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	testStructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

const SuffixKey = "suffix"

func cleanup(t *testing.T, terraformOptions *terraform.Options, tempTestFolder string) {
	terraform.Destroy(t, terraformOptions)
	os.RemoveAll(tempTestFolder)
}

// Test the Terraform module in examples/complete using Terratest.
func TestExamplesComplete(t *testing.T) {
	t.Parallel()

	// Uncomment these when doing local testing if you need to skip any stages.
	//os.Setenv("SKIP_bootstrap", "true")
	//os.Setenv("SKIP_apply", "true")
	//os.Setenv("SKIP_destroy", "true")

	rootFolder := "../../"

	terraformFolderRelativeToRoot := "examples/complete"

	tempTestFolder := testStructure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer testStructure.RunTestStage(t, "teardown", func() {
		terraformOptions := testStructure.LoadTerraformOptions(t, tempTestFolder)
		cleanup(t, terraformOptions, tempTestFolder)
	})

	testStructure.RunTestStage(t, "bootstrap", func() {
		randID := strings.ToLower(random.UniqueId())
		testStructure.SaveString(t, tempTestFolder, SuffixKey, randID)
	})

	// Apply the infrastructure
	testStructure.RunTestStage(t, "apply", func() {
		suffix := testStructure.LoadString(t, tempTestFolder, SuffixKey)

		terraformOptions := &terraform.Options{
			// The path to where our Terraform code is located
			TerraformDir:    tempTestFolder,
			TerraformBinary: "terragrunt",
			Upgrade:         true,
			// Variables to pass to our Terraform code using -var-file options
			Vars: map[string]interface{}{
				"enabled": "true",
				"suffix":  suffix,
			},
		}

		// Save the terraform oprions for future reference
		testStructure.SaveTerraformOptions(t, tempTestFolder, terraformOptions)
		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
		terraform.InitAndApply(t, terraformOptions)

		// Run `terraform output` to get the value of an output variable
		//id := terraform.Output(t, terraformOptions, "vpc_id")
		//name := terraform.Output(t, terraformOptions, "name")
		//arn := terraform.Output(t, terraformOptions, "vpc_arn")
		//
		//// Verify we're getting back the outputs we expect
		//// Ensure we get the attribute included in the ID
		//assert.NotNil(t, id, "example-test")
		//assert.Contains(t, name, "example-test")
		//assert.NotNil(t, arn, "example-test")
		//
		//// This will run `terraform apply` a second time and fail the test if there are any errors
		//terraform.Apply(t, terraformOptions)
		//
		//id2 := terraform.Output(t, terraformOptions, "vpc_id")
		//arn2 := terraform.Output(t, terraformOptions, "vpc_arn")
		//
		//assert.Equal(t, id, id2, "Expected `id` to be stable")
		//assert.Equal(t, arn, arn2, "Expected `arn` to be stable")
	})

	// Run perpetual diff
	testStructure.RunTestStage(t, "perpetual_diff", func() {
		terraformOptions := testStructure.LoadTerraformOptions(t, tempTestFolder)
		planResult := terraform.Plan(t, terraformOptions)

		// Make sure the plan shows zero changes
		assert.Contains(t, planResult, "No changes.")
	})
}

func TestExamplesCompleteDisabled(t *testing.T) {
	t.Parallel()

	rootFolder := "../../"
	terraformFolderRelativeToRoot := "examples/complete"
	// Uncomment these when doing local testing if you need to skip any stages.
	//os.Setenv("SKIP_apply", "true")
	//os.Setenv("SKIP_destroy", "true")

	tempTestFolder := testStructure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	// Apply the infrastructure
	testStructure.RunTestStage(t, "apply", func() {
		terraformOptions := &terraform.Options{
			// The path to where our Terraform code is located
			TerraformDir:    tempTestFolder,
			Upgrade:         true,
			TerraformBinary: "terragrunt",
			Vars: map[string]interface{}{
				"enabled": "false",
			},
		}

		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
		results := terraform.InitAndApply(t, terraformOptions)

		// Should complete successfully without creating or changing any resources
		assert.Contains(t, results, "Resources: 0 added, 0 changed, 0 destroyed.")
	})
}
