package template

func defaultProvidersConfig() map[string]map[string]string {
	return map[string]map[string]string{
		"azurerm": {
			"source":  "hashicorp/azurerm",
			"version": "~>3.0",
		},
		"random": {
			"source":  "hashicorp/random",
			"version": "~>3.0",
		},
	}
}
