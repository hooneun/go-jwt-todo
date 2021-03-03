package auth

type JWTInterface interface {
	CreateToken(uint64) (JWTToken, error)
	SetTokenRedis(uint64, JWTToken) error
}
