resource "azurerm_resource_group" "web_app_mysql_rg" {
  name     = "${random_pet.prefix.id}-rg"
  location = azurerm_resource_group.web_app_mysql_rg.location
  tags     = var.tags
}

resource "random_pet" "prefix" {
  prefix = var.prefix
  length = 1
}
resource "azurerm_mysql_server" "web_app_backend" {
  name                = "${replace(random_pet.prefix.id, "-", "")}pgserver"
  location            = azurerm_resource_group.web_app_mysql_rg.location
  resource_group_name = azurerm_resource_group.web_app_mysql_rg.name
  tags                = azurerm_resource_group.web_app_mysql_rg.tags

  administrator_login          = "${var.prefix}-admin"
  administrator_login_password = random_password.password.result

  sku_name   = var.database_sku_name
  storage_mb = var.database_sku_size_MB
  version    = var.mysql_version

  ssl_enforcement_enabled          = false
  ssl_minimal_tls_version_enforced = "TLSEnforcementDisabled"
}

resource "azurerm_mysql_database" "web_app_backend" {
  name                = "${replace(random_pet.prefix.id, "-", "")}database"
  resource_group_name = azurerm_resource_group.web_app_mysql_rg.name

  server_name = azurerm_mysql_server.web_app_backend.name
  charset     = "utf8mb4"
  collation   = "utf8mb4_unicode_ci"
}

resource "random_password" "password" {
  length      = 20
  min_lower   = 1
  min_upper   = 1
  min_numeric = 1
  min_special = 1
  special     = false
}
resource "azurerm_app_service_plan" "web_app_frontend" {
  name                = "${replace(random_pet.prefix.id, "-", "")}serviceplan"
  resource_group_name = azurerm_resource_group.web_app_mysql_rg.name
  location            = azurerm_resource_group.web_app_mysql_rg.location
  tags                = azurerm_resource_group.web_app_mysql_rg.tags

  sku {
    tier = var.service_plan_tier
    size = var.service_plan_size
  }
}

resource "azurerm_app_service" "main" {
  name                = "${replace(random_pet.prefix.id, "-", "")}service"
  location            = azurerm_resource_group.web_app_mysql_rg.location
  resource_group_name = azurerm_resource_group.web_app_mysql_rg.name
  tags                = azurerm_resource_group.web_app_mysql_rg.tags

  app_service_plan_id = azurerm_app_service_plan.web_app_frontend.id
  connection_string {
    name  = "DefaultConnect"
    type  = "MySql"
    value = "Database=${azurerm_mysql_database.web_app_backend.name};Data Source=${azurerm_mysql_server.web_app_backend.fqdn};User Id=${random_pet.prefix.id}-admin@${azurerm_mysql_server.web_app_backend.name};Password=${random_password.password.result}"
  }
}
