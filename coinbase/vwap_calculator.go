package coinbase

import (
	"context"
	"math/big"

	"github.com/flaviostutz/coinbase-vwap/mathutils"
)

type VWAPInfo struct {
	ProductId string
	Value     *big.Float
}

type OnVWAPInfo func(VWAPInfo)

//CalculateVWAP Consumes a matches channel and for each sample calculates the instantaneous VWAP for each
//ProductId that comes in the channel. For each received sample, the callback function onVMAPInfo will be
//called with the latest results in weighted average for the price, according to VWAP
func CalculateVWAP(ctx context.Context, matchInfoIn <-chan MatchInfo, averagerMaxSize int, onVMAPInfo OnVWAPInfo) error {
	averagers := make(map[string]*mathutils.WeightedAvg)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				//resolve an averager for the received product id
				mi := <-matchInfoIn
				av, ok := averagers[mi.ProductId]
				if !ok {
					a := mathutils.NewWeightedAvg(averagerMaxSize)
					av = &a
					averagers[mi.ProductId] = av
				}

				//add value to averager and calculate current weighted average
				av.Add(mi.Price, mi.Size)
				onVMAPInfo(VWAPInfo{
					ProductId: mi.ProductId,
					Value:     av.Avg(),
				})
			}
		}
	}()

	return nil
}
