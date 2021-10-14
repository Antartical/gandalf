package security

import (
	"crypto/rand"
	"io"
	"math/big"
)

const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-?@#!="

/*
ISecretGenerator -> interface for secrets generator
*/
type ISecretGenerator interface {
	GenerateSecret(lenght int) (string, error)
}

/*
UniformSecretGenerator -> secret based on uniform character selection
*/
type UniformSecretGenerator struct {
	getRandomPosition func(rand io.Reader, max *big.Int) (n *big.Int, err error)
}

/*
GenerateSecret -> returns securely generated random string
*/
func (secret UniformSecretGenerator) GenerateSecret(lenght int) (string, error) {
	ret := make([]byte, lenght)
	for i := 0; i < lenght; i++ {
		num, err := secret.getRandomPosition(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			return "", err
		}
		ret[i] = characters[num.Int64()]
	}

	return string(ret), nil
}

/*
NewUniformSecret -> creates a new uniform secret
*/
func NewUniformSecret() UniformSecretGenerator {
	return UniformSecretGenerator{
		getRandomPosition: rand.Int,
	}
}
