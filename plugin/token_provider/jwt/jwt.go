package jwt

import (
	"flag"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"todo-service/common"
	token_provider "todo-service/plugin/token_provider"
)

type jwtProvider struct {
	name   string
	secret string
}

func NewJWTProvider(name string) *jwtProvider {
	return &jwtProvider{
		name: name,
	}
}

func (p *jwtProvider) GetPrefix() string {
	return p.name
}

func (p *jwtProvider) Get() interface{} {
	return p
}

func (p *jwtProvider) Name() string {
	return p.name
}

func (p *jwtProvider) InitFlags() {
	flag.StringVar(&p.secret, "jwt-secret", "whoant", "Secret key for generating JWT token")
}

func (p *jwtProvider) Configure() error {
	return nil
}

func (p *jwtProvider) Run() error {
	return nil
}

func (p *jwtProvider) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (p *jwtProvider) SecretKey() string {
	return p.secret
}

func (p *jwtProvider) Generate(data token_provider.TokenPayload, expiry int) (token_provider.Token, error) {
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		common.TokenPayload{
			UId:   data.UserId(),
			URole: data.Role(),
		},
		jwt.StandardClaims{
			ExpiresAt: now.Local().Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt:  now.Local().Unix(),
			Id:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})

	myToken, err := t.SignedString([]byte(p.secret))
	if err != nil {
		return nil, err
	}

	return &token{
		Token:   myToken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (token_provider.TokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, token_provider.ErrInvalidToken
	}

	if !res.Valid {
		return nil, token_provider.ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, token_provider.ErrInvalidToken
	}

	return claims.Payload, nil
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.StandardClaims
}

type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}
