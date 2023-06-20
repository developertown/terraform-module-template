//go:build unit
// +build unit

package test

import (
	"github.com/gruntwork-io/terratest/modules/random"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	testStructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// Test the Terraform module in examples/complete using Terratest.
func TestVPCCompleteValidation(t *testing.T) {
	t.Parallel()

	rootFolder := "../../"
	randID := strings.ToLower(random.UniqueId())
	terraformFolderRelativeToRoot := "examples/complete"

	tempTestFolder := testStructure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir:    tempTestFolder,
		TerraformBinary: "terragrunt",
		Upgrade:         true,
		// Variables to pass to our Terraform code using -var-file options
		Vars: map[string]interface{}{
			"enabled": "true",
			"suffix":  randID,
		},
	}

	output := terraform.InitAndPlan(t, terraformOptions)
	assert.Contains(t, output, "50 to add, 0 to change, 0 to destroy", "Plan OK and should attempt to create 50 resources")
}
