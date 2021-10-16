package security

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

type getRandomPositionRecorder struct {
	reader io.Reader
	max    *big.Int
}

func mockedGetRandomPosition(err error) (func(rand io.Reader, max *big.Int) (n *big.Int, err error), *getRandomPositionRecorder) {
	recorder := new(getRandomPositionRecorder)
	return func(rand io.Reader, max *big.Int) (*big.Int, error) {
		*recorder = getRandomPositionRecorder{rand, max}
		return big.NewInt(int64(0)), err
	}, recorder
}

func TestGenerateSecret(t *testing.T) {
	assert := require.New(t)

	t.Run("Test uniform secret constructor", func(t *testing.T) {
		lenght := 5
		uniformGenerator := NewUniformSecret()
		secret, err := uniformGenerator.GenerateSecret(lenght)

		assert.Equal(lenght, len(secret))
		assert.Nil(err)
	})

	t.Run("Test generate secret correct", func(t *testing.T) {
		getRandomPosition, recorder := mockedGetRandomPosition(nil)
		secretGenerator := UniformSecretGenerator{
			getRandomPosition: getRandomPosition,
		}
		lenght := 5
		secret, err := secretGenerator.GenerateSecret(lenght)

		assert.Equal(lenght, len(secret))
		assert.Equal(recorder.max, big.NewInt(int64(len(characters))))
		assert.Nil(err)
	})

	t.Run("Test generate secret error", func(t *testing.T) {
		expectedError := errors.New("woohps")
		getRandomPosition, _ := mockedGetRandomPosition(expectedError)
		secretGenerator := UniformSecretGenerator{
			getRandomPosition: getRandomPosition,
		}
		secret, err := secretGenerator.GenerateSecret(5)
		fmt.Print(secret)

		assert.Equal(expectedError, err)
	})

}
