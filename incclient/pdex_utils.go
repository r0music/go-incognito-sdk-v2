package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"math/big"
	"sort"
	"strings"

	// "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	// "github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// Share represents a pDEX contribution share.
type Share struct {
	TokenID1Str string
	TokenID2Str string
	ShareAmount uint64
}

// GetPdexState retrieves the state of pDEX at the provided beacon height.
// If the beacon height is set to 0, it returns the latest pDEX state.
func (client *IncClient) GetPdexState(beaconHeight uint64, filter map[string]interface{}) (*jsonresult.CurrentPdexState, error) {
	if beaconHeight == 0 {
		bestBlocks, err := client.GetBestBlock()
		if err != nil {
			return nil, fmt.Errorf("cannot get best blocks: %v", err)
		}
		beaconHeight = bestBlocks[-1]
	}

	responseInBytes, err := client.rpcServer.GetPdexState(beaconHeight, filter)
	if err != nil {
		return nil, err
	}

	var pdeState jsonresult.CurrentPdexState
	err = rpchandler.ParseResponse(responseInBytes, &pdeState)
	if err != nil {
		return nil, err
	}

	return &pdeState, nil
}

// GetAllPdexPoolPairs retrieves all pools in pDEX at the provided beacon height.
// If the beacon height is set to 0, it returns the latest pDEX pool pairs.
func (client *IncClient) GetAllPdexPoolPairs(beaconHeight uint64) (map[string]*jsonresult.Pdexv3PoolPairState, error) {
	pdeState, err := client.GetPdexState(beaconHeight, nil)
	if err != nil {
		return nil, err
	}

	return pdeState.PoolPairs, nil
}

// GetPdexPoolPair retrieves the pDEX pool information for pair tokenID1-tokenID2 at the provided beacon height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetPdexPoolPair(beaconHeight uint64, tokenID1, tokenID2 string) (map[string]*jsonresult.Pdexv3PoolPairState, error) {
	if beaconHeight == 0 {
		bestBlocks, err := client.GetBestBlock()
		if err != nil {
			return nil, fmt.Errorf("cannot get best blocks: %v", err)
		}
		beaconHeight = bestBlocks[-1]
	}

	allPoolPairs, err := client.GetAllPdexPoolPairs(beaconHeight)
	if err != nil {
		return nil, err
	}

	results := make(map[string]*jsonresult.Pdexv3PoolPairState)
	// search for the pool by concatenating tokenIDs. Try both orderings
	prefix1 := fmt.Sprintf("%s-%s", tokenID1, tokenID2)
	prefix2 := fmt.Sprintf("%s-%s", tokenID2, tokenID1)
	for k, v := range allPoolPairs {
		if strings.HasPrefix(k, prefix1) || strings.HasPrefix(k, prefix2) {
			results[k] = v
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("cannot found pool pair for tokenID %v and %v", tokenID1, tokenID2)
	}
	return results, nil
}

func calculateBuyAmount(amountIn uint64, virtualReserveIn *big.Int, virtualReserveOut *big.Int) (uint64, error) {
	if amountIn <= 0 {
		return 0, fmt.Errorf("invalid input amount %d", amountIn)
	}
	amount := big.NewInt(0).SetUint64(amountIn)
	num := big.NewInt(0).Mul(amount, virtualReserveOut)
	den := big.NewInt(0).Add(amount, virtualReserveIn)
	result := num.Div(num, den)
	if !result.IsUint64() {
		return 0, fmt.Errorf("buy amount out %s of uint64 range", result.String())
	}
	return result.Uint64(), nil
}

// CheckPrice gets the remote server to check price for trading things.
func (client *IncClient) CheckPrice(pairID, tokenToSell string, sellAmount uint64) (uint64, error) {
	pairs, err := client.GetAllPdexPoolPairs(0)
	if err != nil {
		return 0, err
	}
	pair, exists := pairs[pairID]
	if !exists {
		return 0, fmt.Errorf("No pool found for ID %s", pairID)
	}

	var virtualAmtSell, virtualAmtBuy *big.Int
	switch tokenToSell {
	case pair.State.Token0ID.String():
		virtualAmtSell = big.NewInt(0).Set(pair.State.Token0VirtualAmount)
		virtualAmtBuy = big.NewInt(0).Set(pair.State.Token1VirtualAmount)
	case pair.State.Token1ID.String():
		virtualAmtSell = big.NewInt(0).Set(pair.State.Token1VirtualAmount)
		virtualAmtBuy = big.NewInt(0).Set(pair.State.Token0VirtualAmount)
	default:
		return 0, fmt.Errorf("No tokenID %s in pool %s", tokenToSell, pairID)
	}

	buyAmount, err := calculateBuyAmount(sellAmount, virtualAmtSell, virtualAmtBuy)
	if err != nil {
		return 0, err
	}
	return buyAmount, nil
}

// CheckNFTMintingStatus retrieves the status of a (pDEX) NFT minting transaction.
func (client *IncClient) CheckNFTMintingStatus(txHash string) (bool, string, error) {
	responseInBytes, err := client.rpcServer.CheckNFTMintingStatus(txHash)
	if err != nil {
		return false, "", err
	}
	type TmpResult struct {
		ID string `json:"NftID"`
		Status int `json:"Status"`
	}
	var res TmpResult
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return false, "", err
	}
	if res.Status != 1 {
		return false, "", fmt.Errorf("minting failed with status %v", res.Status)
	}

	return true, res.ID, nil
}

// CheckTradeStatus checks the status of a trading transaction.
// It returns
//	- -1: if an error occurred;
//	- 1: if the trade is accepted;
//	- 2: if the trade is not accepted.
func (client *IncClient) CheckTradeStatus(txHash string) (int, error) {
	responseInBytes, err := client.rpcServer.CheckTradeStatus(txHash)
	if err != nil {
		return -1, err
	}

	var tradeStatus rpc.DEXTradeStatus
	err = rpchandler.ParseResponse(responseInBytes, &tradeStatus)
	if err != nil {
		return -1, err
	}

	return tradeStatus.Status, nil
}

// BuildPdexShareKey constructs a key for retrieving contributed shares in pDEX.
func BuildPdexShareKey(beaconHeight uint64, token1ID string, token2ID string, contributorAddress string) ([]byte, error) {
	pdeSharePrefix := []byte("pdeshare-")
	prefix := append(pdeSharePrefix, []byte(fmt.Sprintf("%d-", beaconHeight))...)
	tokenIDs := []string{token1ID, token2ID}
	sort.Strings(tokenIDs)

	var keyAddr string
	var err error
	if len(contributorAddress) == 0 {
		keyAddr = contributorAddress
	} else {
		//Always parse the contributor address into the oldest version for compatibility
		keyAddr, err = wallet.GetPaymentAddressV1(contributorAddress, false)
		if err != nil {
			return nil, err
		}
	}
	return append(prefix, []byte(tokenIDs[0]+"-"+tokenIDs[1]+"-"+keyAddr)...), nil
}

// BuildPdexPoolKey constructs a key for a pool in pDEX.
func BuildPdexPoolKey(token1ID string, token2ID string) string {
	tokenIDs := []string{token1ID, token2ID}
	sort.Strings(tokenIDs)

	return fmt.Sprintf("%v-%v", tokenIDs[0], tokenIDs[1])
}
