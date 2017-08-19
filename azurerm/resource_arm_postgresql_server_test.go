package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPostgreSQLServer_basicNinePointFive(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointFive(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicNinePointSix(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updatePassword(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri)
	updatedConfig := testAccAzureRMPostgreSQLServer_basicNinePointSixUpdated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLServerExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Server: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).postgresqlServersClient

		resp, err := client.Get(resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on postgresqlServersClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: PostgreSQL Server %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).postgresqlServersClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(resourceGroup, name)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("PostgreSQL Server still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMPostgreSQLServer_basicNinePointFive(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}
resource "azurerm_postgresql_server" "test" {
  name = "acctestpsqlsvr-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "PGSQLB50"
    capacity = 50
    tier = "Basic"
  }

  administrator_login = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "9.5"
  storage_mb = 51200
  ssl_enforcement = "Enabled"
}
`, rInt, rInt)
}

func testAccAzureRMPostgreSQLServer_basicNinePointSix(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}
resource "azurerm_postgresql_server" "test" {
  name = "acctestpsqlsvr-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "PGSQLB50"
    capacity = 50
    tier = "Basic"
  }

  administrator_login = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "9.6"
  storage_mb = 51200
  ssl_enforcement = "Enabled"
}
`, rInt, rInt)
}

func testAccAzureRMPostgreSQLServer_basicNinePointSixUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}
resource "azurerm_postgresql_server" "test" {
  name = "acctestpsqlsvr-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "PGSQLB50"
    capacity = 50
    tier = "Basic"
  }

  administrator_login = "acctestun"
  administrator_login_password = "R3dH0TCh1l1P3pp3rs!"
  version = "9.6"
  storage_mb = 51200
  ssl_enforcement = "Disabled"
}
`, rInt, rInt)
}