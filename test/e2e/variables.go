//go:build e2e

package e2e

import (
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var (
	// TestImageName is the name of the system image that is used in end-to-end tests.
	TestImageName = getEnv("TEST_IMAGE", "ubuntu-24.04")

	// TestImageID is the system image ID that is used in end-to-end tests.
	TestImageID = getEnv("TEST_IMAGE_ID", "161547269")

	// TestServerType is the default server type used in end-to-end tests.
	TestServerType = getEnv("TEST_SERVER_TYPE", "cpx22")

	// TestServerTypeUpgrade is the upgrade server type used in end-to-end tests.
	TestServerTypeUpgrade = getEnv("TEST_SERVER_TYPE_UPGRADE", "cpx32")

	// TestArchitecture is the default architecture used in end-to-end tests, should match the architecture of the TestServerType.
	TestArchitecture = getEnv("TEST_ARCHITECTURE", string(hcloud.ArchitectureX86))

	// TestLoadBalancerTypeName is the default Load Balancer type used in end-to-end tests.
	TestLoadBalancerTypeName = getEnv("TEST_LOAD_BALANCER_TYPE_NAME", "lb11")

	// TestDatacenterName is the default datacenter name where we execute our end-to-end tests.
	TestDatacenterName = getEnv("TEST_DATACENTER_NAME", "nbg1-dc3")

	// TestDatacenterID is the default datacenter ID where we execute our end-to-end tests (Must be the ID of TestDatacenterName)
	TestDatacenterID = getEnv("TEST_DATACENTER_ID", "2")

	// TestLocationName is the default location where we execute our end-to-end tests.
	TestLocationName = getEnv("TEST_LOCATION", "nbg1")

	// TestISOName is the default ISO used for testing
	TestISOName = getEnv("TEST_ISO_NAME", "SW_DVD9_Win_Server_STD_CORE_2022__64Bit_German_DC_STD_MLF_X22-74292.iso")

	// TestNetworkZone is the default network zone used for testing
	TestNetworkZone = getEnv("TEST_NETWORK_ZONE", "eu-central")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
