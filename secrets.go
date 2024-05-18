package cleanconfig

type SecretManager interface {
	GetSecret(secret string) (string, error)
}
