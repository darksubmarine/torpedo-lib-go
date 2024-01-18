package conf

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"os"
	"testing"
)

func TestAZAppConfigLoader_Load(t *testing.T) {
	t.Skip() // only run it manually
	// AppConfig connection string
	endpoint := os.Getenv("AZURE_APPCONFIGURATION_ENDPOINT")

	// azIdentity credentials... needed to connect to Azure Key Vault
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("failed to obtain a credential: %v", err)
	}

	loader := NewAZAppConfigLoader("*", "*", endpoint, cred)

	m := Map{}
	m = loader.Load(m)

	fmt.Println(m)
}
