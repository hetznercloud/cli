//go:build e2e

package e2e

import (
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var (
	// TestImage is the system image that is used in end-to-end tests.
	TestImage = getEnv("TEST_IMAGE", "ubuntu-24.04")

	// TestImageID is the system image ID that is used in end-to-end tests.
	TestImageID = getEnv("TEST_IMAGE_ID", "161547269")

	// TestServerType is the default server type used in end-to-end tests.
	TestServerType = getEnv("TEST_SERVER_TYPE", "cpx11")

	// TestServerTypeUpgrade is the upgrade server type used in end-to-end tests.
	TestServerTypeUpgrade = getEnv("TEST_SERVER_TYPE_UPGRADE", "cpx21")

	// TestArchitecture is the default architecture used in end-to-end tests, should match the architecture of the TestServerType.
	TestArchitecture = getEnv("TEST_ARCHITECTURE", string(hcloud.ArchitectureX86))

	// TestLoadBalancerType is the default Load Balancer type used in end-to-end tests.
	TestLoadBalancerType = "lb11"

	// TestDataCenter is the default datacenter where we execute our end-to-end tests.
	TestDataCenter = getEnv("TEST_DATACENTER", "nbg1-dc3")

	// TestLocationName is the default location where we execute our end-to-end tests.
	TestLocationName = getEnv("TEST_LOCATION", "nbg1")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
