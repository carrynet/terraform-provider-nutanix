package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixVirtualMachine_basic(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "staging"),
				),
			},
			{
				Config: testAccNutanixVMConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "production"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_list"},
			},
		},
	})
}

func TestAccNutanixVirtualMachine_WithDisk(t *testing.T) {
	r := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithDisk(r),
			},
			{
				Config:             testAccNutanixVMConfigWithDisk(r),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccNutanixVMConfigWithDiskUpdate(r),
			},
			{
				Config:             testAccNutanixVMConfigWithDiskUpdate(r),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:            "nutanix_virtual_machine.vm1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_list"},
			},
		}})
}

func TestAccNutanixVirtualMachine_updateFields(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm2"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigUpdatedFields(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-dou-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "staging"),
				),
			},
			{
				Config: testAccNutanixVMConfigUpdatedFieldsUpdated(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-dou-%d-updated", r)),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "256"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "production"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_list"},
			},
		},
	})
}

func TestAccNutanixVirtualMachine_wiithSubnet(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithSubnet(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "staging"),
					resource.TestCheckResourceAttrSet(resourceName, "nic_list.0.ip_endpoint_list.0.ip"),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine" {
			continue
		}
		for {
			_, err := conn.API.V3.GetVM(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}

	}

	return nil

}

func testAccNutanixVMConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}
output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-%d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186


	categories {
		environment-terraform = "staging"
	}
}
`, r)
}

func testAccNutanixVMConfigWithDisk(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_image" "cirros-034-disk" {
    name        = "test-image-dou-%[1]d"
    source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
    description = "heres a tiny linux image, not an iso, but a real disk!"
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-vm-%[1]d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186

	disk_list = [{
		# data_source_reference in the Nutanix API refers to where the source for
		# the disk device will come from. Could be a clone of a different VM or a
		# image like we're doing here
		data_source_reference = [{
			kind = "image"
			uuid = "${nutanix_image.cirros-034-disk.id}"
		}]
		disk_size_mib = 44
	},
	{
		disk_size_mib = 100
	},
	{
		disk_size_mib = 200
	},
	{
		disk_size_mib = 300
	}]
}
`, r)
}

func testAccNutanixVMConfigWithDiskUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_image" "cirros-034-disk" {
    name        = "test-image-dou-%[1]d"
    source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
    description = "heres a tiny linux image, not an iso, but a real disk!"
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-vm-%[1]d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186

	disk_list = [{
		# data_source_reference in the Nutanix API refers to where the source for
		# the disk device will come from. Could be a clone of a different VM or a
		# image like we're doing here
		data_source_reference = [{
			kind = "image"
			uuid = "${nutanix_image.cirros-034-disk.id}"
		}]
		disk_size_mib = 44
	},
	{
		disk_size_mib = 100
	},
	{
		disk_size_mib = 400
	}]
}
`, r)
}

func testAccNutanixVMConfigUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}


output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-%d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 2
  memory_size_mib      = 186

	categories {
		environment-terraform = "production"
	}
}
`, r)
}

func testAccNutanixVMConfigUpdatedFields(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}
output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_virtual_machine" "vm2" {
  name = "test-dou-%d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186


	categories {
		environment-terraform = "staging"
	}
}
`, r)
}

func testAccNutanixVMConfigUpdatedFieldsUpdated(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}


output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_virtual_machine" "vm2" {
  name = "test-dou-%d-updated"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 2
  num_sockets          = 2
  memory_size_mib      = 256

	categories {
		environment-terraform = "production"
	}
}
`, r)
}

func testAccNutanixVMConfigWithSubnet(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

resource "nutanix_subnet" "sub" {
  cluster_reference = {
		kind = "cluster"
		uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed-%[1]d"
	description = "Description of my unit test VLAN"
  vlan_id     = 50
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
  ip_config_pool_list_ranges   = ["10.250.140.20 10.250.140.100"]
}

resource "nutanix_image" "cirros-034-disk" {
    name        = "test-image-dou-%[1]d"
    source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
    description = "heres a tiny linux image, not an iso, but a real disk!"
}

resource "nutanix_virtual_machine" "vm1" {
	name = "test-dou-vm-%[1]d"
	
	categories {
		environment-terraform = "staging"
	}

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186

	disk_list = [{
		data_source_reference = [{
			kind = "image"
			uuid = "${nutanix_image.cirros-034-disk.id}"
		}]
		disk_size_mib = 44
	}]

	nic_list = [{
		subnet_reference = {
			kind = "subnet"
			uuid = "${nutanix_subnet.sub.id}"
		}
	}]
}
`, r)
}
