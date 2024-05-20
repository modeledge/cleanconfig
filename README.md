# cleanconfig

## License

This project is based on [Original Project Name](link-to-original-project) and is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more information.

## Description

Fork of https://github.com/ilyakaznacheev/cleanenv with an optional secret manager which implements the SecretManger interface:

```go
type SecretManager interface {
	GetSecret(secret string) (string, error)
}
```
