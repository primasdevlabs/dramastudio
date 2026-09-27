package security

type TokenValidator interface {
	ValidateToken(token string) (string, error)
}
