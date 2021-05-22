package mathutils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedAvg1(t *testing.T) {
	wa := NewWeightedAvg(1)
	wa.Add(big.NewFloat(1), big.NewFloat(1))
	assert.Equal(t, "1", wa.Avg().String())

	wa.Add(big.NewFloat(2), big.NewFloat(2))
	assert.Equal(t, "2", wa.Avg().String())

	wa.Add(big.NewFloat(3), big.NewFloat(3))
	assert.Equal(t, "3", wa.Avg().String())
}

func TestWeightedAvg31(t *testing.T) {
	wa := NewWeightedAvg(3)
	wa.Add(big.NewFloat(1), big.NewFloat(1))
	assert.Equal(t, "1", wa.Avg().String())

	wa.Add(big.NewFloat(3), big.NewFloat(1))
	assert.Equal(t, "2", wa.Avg().String())

	wa.Add(big.NewFloat(2), big.NewFloat(1))
	assert.Equal(t, "2", wa.Avg().String())

	wa.Add(big.NewFloat(7), big.NewFloat(1))
	assert.Equal(t, "4", wa.Avg().String())

	wa.Add(big.NewFloat(15), big.NewFloat(1))
	assert.Equal(t, "8", wa.Avg().String())
}

func TestWeightedAvg123(t *testing.T) {
	wa := NewWeightedAvg(2)
	wa.Add(big.NewFloat(1), big.NewFloat(1))
	assert.Equal(t, "1", wa.Avg().String())

	wa.Add(big.NewFloat(2), big.NewFloat(2))
	assert.Equal(t, "1.666666667", wa.Avg().String())

	wa.Add(big.NewFloat(3), big.NewFloat(3))
	assert.Equal(t, "2.6", wa.Avg().String())
}

func TestWeightedAvg4(t *testing.T) {
	wa := NewWeightedAvg(4)

	wa.Add(big.NewFloat(1), big.NewFloat(2))
	assert.Equal(t, "1", wa.Avg().String())
	wa.Add(big.NewFloat(1), big.NewFloat(4))
	assert.Equal(t, "1", wa.Avg().String())
	wa.Add(big.NewFloat(3), big.NewFloat(2))
	assert.Equal(t, "1.5", wa.Avg().String())
	wa.Add(big.NewFloat(2), big.NewFloat(6))
	assert.Equal(t, "1.714285714", wa.Avg().String())
	//this should remove the first sample from calc
	wa.Add(big.NewFloat(3), big.NewFloat(1))
	assert.Equal(t, "1.923076923", wa.Avg().String())
}
