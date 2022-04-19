package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"agile": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

// This provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// PreCheck(t) must be called before using this provider instance.
var testAccProvider *schema.Provider = New("dev")()

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))

	if err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("AGILE_USERNAME"); v == "" {
		t.Fatal("AGILE_USERNAME env variable must be set for acceptance tests")
	}

	if v := os.Getenv("AGILE_PASSWORD"); v == "" {
		t.Fatal("AGILE_PASSWORD env variable must be set for acceptance tests")
	}

	if v := os.Getenv("AGILE_API"); v == "" {
		t.Fatal("AGILE_API env variable must be set for acceptance tests")
	}
}
