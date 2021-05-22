package coinbase

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestVWAPCalculator1(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

	//invoke calculation method and verify results
	mic := make(chan MatchInfo)
	counter := 1
	CalculateVWAP(ctx, mic, 2, func(info VWAPInfo) {
		switch {
		case counter == 1:
			assert.Equal(t, "PRD1", info.ProductId)
			assert.Equal(t, "1", info.Value.String())
		case counter == 2:
			assert.Equal(t, "PRD1", info.ProductId)
			assert.Equal(t, "1.666666667", info.Value.String())
		case counter == 3:
			assert.Equal(t, "PRD1", info.ProductId)
			assert.Equal(t, "2.6", info.Value.String())
			//reached end of test. finish context (won't wait for timeout)
			cancel()
		default:
			assert.Fail(t, "Should not reach here")
		}
		counter++
	})

	//produce 3 messages to chan
	mic <- MatchInfo{
		ProductId: "PRD1",
		Price:     big.NewFloat(1),
		Size:      big.NewFloat(1),
	}
	mic <- MatchInfo{
		ProductId: "PRD1",
		Price:     big.NewFloat(2),
		Size:      big.NewFloat(2),
	}
	mic <- MatchInfo{
		ProductId: "PRD1",
		Price:     big.NewFloat(3),
		Size:      big.NewFloat(3),
	}

	select {
	case <-ctx.Done():
		//wait for context timeout/cancelation
	}
}
