package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"net/url"
	"strings"
)

type AZAppConfigLoader struct {
	endpoint         string
	credentials      *azidentity.DefaultAzureCredential
	keyFilter        string
	labelFilter      string
	azsecretsClients map[string]*azsecrets.Client
}

func NewAZAppConfigLoader(keyFilter string, labelFilter string, conn string, cred *azidentity.DefaultAzureCredential) *AZAppConfigLoader {
	return &AZAppConfigLoader{endpoint: conn, credentials: cred, keyFilter: keyFilter, labelFilter: labelFilter, azsecretsClients: map[string]*azsecrets.Client{}}
}

func (l *AZAppConfigLoader) Load(m Map) Map {

	settings, err := l.AZAppConfigListSettings()
	if err != nil {
		return m
	}

	for k, v := range settings {
		m.Add(v, strings.Split(strings.Trim(k, "/"), "/")[1:]...)
	}

	return m
}

func (l *AZAppConfigLoader) AZAppConfigListSettings() (map[string]string, error) {
	client, err := azappconfig.NewClient(l.endpoint, l.credentials, nil)
	if err != nil {
		return nil, err
	}

	revPgr := client.NewListRevisionsPager(
		azappconfig.SettingSelector{
			KeyFilter:   to.Ptr(l.keyFilter),
			LabelFilter: to.Ptr(l.labelFilter),
			Fields:      azappconfig.AllSettingFields(),
		},
		nil)

	settings := map[string]string{}
	for revPgr.More() {
		if revResp, revErr := revPgr.NextPage(context.TODO()); revErr == nil {
			for _, setting := range revResp.Settings {
				if setting.Key != nil && setting.Value != nil {
					if setting.ContentType != nil {
						if val, err := l.AZAppConfigReadSetting(*setting.Value, *setting.ContentType); err != nil {
							//return nil, err
							continue
						} else {
							settings[*setting.Key] = val
						}
					} else {
						settings[*setting.Key] = *setting.Value
					}
				}
			}
		} else {
			return settings, revErr
		}
	}

	return settings, nil
}

func (l *AZAppConfigLoader) AZAppConfigReadSetting(value string, contentType string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("empty value")
	}

	if strings.Contains(contentType, "vnd.microsoft.appconfig.keyvaultref+json") {
		if uri, err := AZAppConfigFetchKeyVaultRefUri(value); err != nil {
			return "", err
		} else {
			if vaultURI, secret, err := AZParseSecretURI(uri); err != nil {
				return "", err
			} else {
				if secretValue, err := l.AZReadSecret(vaultURI, secret); err != nil {
					return "", err
				} else {
					return secretValue, nil
				}
			}
		}
	}

	return value, nil
}

func (l *AZAppConfigLoader) azsecretsClientFor(vaultURI string) (*azsecrets.Client, error) {
	if c, exists := l.azsecretsClients[vaultURI]; exists {
		return c, nil
	}

	client, err := azsecrets.NewClient(vaultURI, l.credentials, &azsecrets.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create a Azure Key Vault client: %w", err)
	}
	l.azsecretsClients[vaultURI] = client
	return client, nil

}

// AZReadSecret reads an Azure Key vault secret value given the vaultURI and the secret name.
func (l *AZAppConfigLoader) AZReadSecret(vaultURI string, secret string) (string, error) {
	// Establish a connection to the Key Vault client
	client, err := l.azsecretsClientFor(vaultURI)
	if err != nil {
		return "", err
	}

	// Get a secret. An empty string version gets the latest version of the secret.
	version := ""
	resp, err := client.GetSecret(context.TODO(), secret, version, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get the secret: %v", err)
	}

	if resp.Value != nil {
		return *resp.Value, nil
	}

	return "", fmt.Errorf("empty value")
}

//----------------------------------
//	Helper functions/structs
//----------------------------------

type KeyVaultRef struct {
	Uri string `json:"uri"`
}

func AZAppConfigFetchKeyVaultRefUri(value string) (string, error) {
	ref := KeyVaultRef{}
	if err := json.Unmarshal([]byte(value), &ref); err != nil {
		return "", err
	} else {
		return ref.Uri, nil
	}
}

// AZParseSecretURI given a secretURI returns its vaultURI (string) and secret name (string) also an error
func AZParseSecretURI(secretURI string) (string, string, error) {
	//uri := "https://your-vault-name.vault.azure.net/secrets/demo"
	if parts, err := url.Parse(secretURI); err != nil {
		return "", "", err
	} else {
		vaultURI := fmt.Sprintf("%s://%s", parts.Scheme, parts.Hostname())
		attrs := strings.Split(parts.Path, "/")
		if len(attrs) == 3 && attrs[1] == "secrets" {
			return vaultURI, attrs[2], nil
		}
		return "", "", fmt.Errorf("secret name cannot be parsed")
	}
}
