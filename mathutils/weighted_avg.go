package mathutils

import (
	"math/big"
)

//WeightedAvg Instantiate with a max number of samples and then by adding samples with
//value + weight you can get the instantaneous weighted average in a rolling window
//Uses big package so it is safe to be used in financial applications
//For performance reasons, we chose to keep partial sums for both SUM(V*W) and SUM(W)
//by removing "expired" elements from its sum while adding new elements to the sum when adding
//a new sample to the window. This strategy reduces the "regular" complexity from O(N) to O(1).
type WeightedAvg struct {
	weightedValues            []*big.Float
	weights                   []*big.Float
	size                      int
	factorTotalWeightedValues *big.Float
	factorTotalWeights        *big.Float
}

func NewWeightedAvg(size int) WeightedAvg {
	w := WeightedAvg{
		weightedValues:            make([]*big.Float, 0),
		weights:                   make([]*big.Float, 0),
		size:                      size,
		factorTotalWeightedValues: new(big.Float),
		factorTotalWeights:        new(big.Float),
	}
	return w
}

func (wa *WeightedAvg) Add(value *big.Float, weight *big.Float) {
	removedWeightedValue := new(big.Float)
	removedWeight := new(big.Float)

	//limit array size
	if len(wa.weightedValues) == wa.size {
		removedWeightedValue = wa.weightedValues[0]
		wa.weightedValues = wa.weightedValues[1:]

		removedWeight = wa.weights[0]
		wa.weights = wa.weights[1:]
	}

	//append values
	mr := new(big.Float).Mul(value, weight)
	wa.weightedValues = append(wa.weightedValues, mr)
	wa.weights = append(wa.weights, weight)

	//pre calculated factors are used for optimization purposes
	//remove outgoing elements and add the newest elements to the total
	ftwv := new(big.Float)
	ftwv.Sub(wa.factorTotalWeightedValues, removedWeightedValue)
	ftwv.Add(ftwv, mr)
	wa.factorTotalWeightedValues = ftwv

	ftw := new(big.Float)
	ftw.Sub(wa.factorTotalWeights, removedWeight)
	ftw.Add(ftw, weight)
	wa.factorTotalWeights = ftw
}

func (wa *WeightedAvg) Avg() *big.Float {
	return new(big.Float).Quo(wa.factorTotalWeightedValues, wa.factorTotalWeights)
}
