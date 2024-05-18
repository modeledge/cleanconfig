package simple_secret

import "fmt"

type FixtureSecrets struct {
	secrets map[string]string
}

func (sm *FixtureSecrets) GetSecret(secretKey string) (string, error) {
	value, ok := sm.secrets[secretKey]
	if !ok {
		return "", fmt.Errorf("secret not found for key: %s", secretKey)
	}
	return value, nil
}

func (sm *FixtureSecrets) SetSecret(secretKey, secretValue string) {
	sm.secrets[secretKey] = secretValue
}

func NewFixtureSecrets() *FixtureSecrets {
	return &FixtureSecrets{
		secrets: map[string]string{
			"MY_SECRET": "ABC",
		},
	}
}
