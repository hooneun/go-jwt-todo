package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type JWTToken struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpire     int64
}

type JWTHandler struct {
	h JWTInterface
}

func NewJWTHandler() (*JWTHandler, error) {
	return new(JWTHandler), nil
}

func (h *JWTHandler) CreateToken(userID uint64) (JWTToken, error) {
	var t JWTToken
	t.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	t.AccessUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = t.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	var err error
	t.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return t, err
	}

	t.RtExpire = time.Now().Add(time.Hour * 24 * 7).Unix()
	t.RefreshUuid = uuid.NewV4().String()

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = t.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = t.RtExpire
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_REFRESH")))
	if err != nil {
		return t, err
	}

	return t, nil
}
