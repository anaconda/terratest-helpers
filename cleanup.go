package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// Cleanup removes the test data folder and the provider.tf if it exits.
func Cleanup(t *testing.T, terraformDir string) {
	test_structure.CleanupTestDataFolder(t, terraformDir)

	// Clean up the test-provider.tf
	providerPath := filepath.Join(terraformDir, "test-provider.tf")
	if files.FileExists(providerPath) {
		os.Remove(providerPath)
	}
}
