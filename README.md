# cleanconfig

Fork of https://github.com/ilyakaznacheev/cleanenv with an optional secret manager which implements the SecretManger interface:

```go
type SecretManager interface {
	GetSecret(secret string) (string, error)
}
```
