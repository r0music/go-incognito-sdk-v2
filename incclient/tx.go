package incclient

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver1"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreateRawTransaction creates a PRV transaction with the provided version.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTransaction(param *TxParam, version int8) ([]byte, string, error) {
	if version == -1 { //Try either one of the version, if possible
		encodedTx, txHash, err := client.CreateRawTransactionVer1(param)
		if err != nil {
			encodedTx, txHash, err1 := client.CreateRawTransactionVer2(param)
			if err1 != nil {
				return nil, "", fmt.Errorf("cannot create raw transaction for either version: %v, %v", err, err1)
			}

			return encodedTx, txHash, nil
		}

		return encodedTx, txHash, nil
	} else if version == 2 {
		return client.CreateRawTransactionVer2(param)
	} else if version == 1 {
		return client.CreateRawTransactionVer1(param)
	}

	return nil, "", fmt.Errorf("transaction version is invalid")
}

// CreateRawTransactionVer1 creates a PRV transaction version 1.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTransactionVer1(param *TxParam) ([]byte, string, error) {
	privateKey := param.senderPrivateKey
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	//Create list of payment infos
	paymentInfos, err := createPaymentInfos(param.receiverList, param.amountList)
	if err != nil {
		return nil, "", err
	}

	//Calculate the total transacted amount
	if param.fee == 0 {
		param.fee = DefaultPRVFee
	}
	totalAmount := param.fee
	for _, amount := range param.amountList {
		totalAmount += amount
	}

	hasPrivacy := true
	if param.md != nil {
		hasPrivacy = false
	}

	coinsToSpend, kvArgs, err := client.initParams(privateKey, common.PRVIDStr, totalAmount, hasPrivacy, 1)
	if err != nil {
		return nil, "", err
	}

	txInitParam := tx_generic.NewTxPrivacyInitParams(&(senderWallet.KeySet.PrivateKey), paymentInfos, coinsToSpend, param.fee, hasPrivacy, &common.PRVCoinID, param.md, nil, kvArgs)

	tx := new(tx_ver1.Tx)
	err = tx.Init(txInitParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txver1 error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txver1: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateRawTransactionVer2 creates a PRV transaction version 2.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTransactionVer2(param *TxParam) ([]byte, string, error) {
	privateKey := param.senderPrivateKey
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	//Create list of payment infos
	paymentInfos, err := createPaymentInfos(param.receiverList, param.amountList)
	if err != nil {
		return nil, "", err
	}

	txFee := param.fee
	if param.fee == 0 {
		txFee = DefaultPRVFee
	}

	//Calculate the total transacted amount
	totalAmount := txFee
	for _, amount := range param.amountList {
		totalAmount += amount
	}

	hasPrivacy := true
	if param.md != nil {
		hasPrivacy = false
	}

	coinsToSpend, kArgs, err := client.initParams(privateKey, common.PRVIDStr, totalAmount, hasPrivacy, 2)
	if err != nil {
		return nil, "", err
	}
	txParam := tx_generic.NewTxPrivacyInitParams(&(senderWallet.KeySet.PrivateKey), paymentInfos, coinsToSpend, DefaultPRVFee, hasPrivacy, &common.PRVCoinID, param.md, nil, kArgs)

	tx := new(tx_ver2.Tx)
	err = tx.Init(txParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txver2 error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txver2: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateAndSendRawTransaction creates a PRV transaction with the provided version, and submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendRawTransaction(privateKey string, addrList []string, amountList []uint64, version int8, md metadata.Metadata) (string, error) {
	txParam := NewTxParam(privateKey, addrList, amountList, 0, nil, md, nil)
	encodedTx, txHash, err := client.CreateRawTransaction(txParam, version)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateRawConversionTransaction creates a PRV transaction that converts PRV coins version 1 to version 2.
// This type of transactions is non-private by default.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawConversionTransaction(privateKey string) ([]byte, string, error) {
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	//Get list of UTXOs
	utxoList, _, err := client.GetUnspentOutputCoins(privateKey, common.PRVIDStr, 0)
	if err != nil {
		return nil, "", err
	}

	//Get list of coinV1 to convert.
	coinV1List, _, _, err := divideCoins(utxoList, nil, true)
	if err != nil {
		return nil, "", fmt.Errorf("cannot divide coin: %v", err)
	}

	if len(coinV1List) == 0 {
		return nil, "", fmt.Errorf("no CoinV1 left to be converted")
	}

	//Calculating the total amount being converted.
	totalAmount := uint64(0)
	for _, utxo := range coinV1List {
		totalAmount += utxo.GetValue()
	}
	if totalAmount < DefaultPRVFee {
		fmt.Printf("Total amount (%v) is less than txFee (%v).\n", totalAmount, DefaultPRVFee)
		return nil, "", fmt.Errorf("Total amount (%v) is less than txFee (%v).\n", totalAmount, DefaultPRVFee)
	}
	totalAmount -= DefaultPRVFee

	uniquePayment := key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Amount: totalAmount, Message: []byte{}}

	//Create tx conversion params
	txParam := tx_ver2.NewTxConvertVer1ToVer2InitParams(&(senderWallet.KeySet.PrivateKey), []*key.PaymentInfo{&uniquePayment}, coinV1List,
		DefaultPRVFee, nil, nil, nil, nil)

	tx := new(tx_ver2.Tx)
	err = tx_ver2.InitConversion(tx, txParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txconvert error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txconvert: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateAndSendRawConversionTransaction creates a PRV transaction that converts PRV coins version 1 to version 2 and broadcasts it to the network.
// This type of transactions is non-private by default.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendRawConversionTransaction(privateKey string, tokenID string) (string, error) {
	var txHash string
	var err error
	var encodedTx []byte

	if tokenID == common.PRVIDStr {
		encodedTx, txHash, err = client.CreateRawConversionTransaction(privateKey)
		if err != nil {
			return "", err
		}

		err = client.SendRawTx(encodedTx)
		if err != nil {
			return "", err
		}
	} else {
		encodedTx, txHash, err = client.CreateRawTokenConversionTransaction(privateKey, tokenID)
		if err != nil {
			return "", err
		}

		err = client.SendRawTokenTx(encodedTx)
		if err != nil {
			return "", err
		}
	}

	return txHash, nil
}

// CreateRawTransactionWithInputCoins creates a raw PRV transaction from the provided input coins.
func (client *IncClient) CreateRawTransactionWithInputCoins(param *TxParam, inputCoins []coin.PlainCoin, coinIndices []uint64) ([]byte, string, error) {
	var txHash string
	senderWallet, err := wallet.Base58CheckDeserialize(param.senderPrivateKey)
	if err != nil {
		return nil, txHash, fmt.Errorf("cannot init private key %v: %v", param.senderPrivateKey, err)
	}
	//Create list of payment infos
	paymentInfos, err := createPaymentInfos(param.receiverList, param.amountList)
	if err != nil {
		return nil, txHash, err
	}
	//Get tx fee
	txFee := DefaultPRVFee
	if param.fee != 0 {
		txFee = param.fee
	}
	//Calculate the total transacted amount
	totalAmount := txFee
	for _, amount := range param.amountList {
		totalAmount += amount
	}
	var kvArgs = make(map[string]interface{})
	if coinIndices == nil {
		//Retrieve commitments and indices
		kvArgs, err = client.getRandomCommitmentV1(inputCoins, common.PRVIDStr)
		if err != nil {
			return nil, txHash, err
		}
		txInitParam := tx_generic.NewTxPrivacyInitParams(&(senderWallet.KeySet.PrivateKey), paymentInfos, inputCoins, txFee, true, &common.PRVCoinID, nil, nil, kvArgs)
		tx := new(tx_ver1.Tx)
		err = tx.Init(txInitParam)
		if err != nil {
			return nil, txHash, fmt.Errorf("init txver1 error: %v", err)
		}
		txBytes, err := json.Marshal(tx)
		if err != nil {
			return nil, txHash, fmt.Errorf("cannot marshal txver1: %v", err)
		}
		base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)
		return []byte(base58CheckData), tx.Hash().String(), nil
	} else {
		//Retrieve commitments and indices
		shardID := GetShardIDFromPrivateKey(param.senderPrivateKey)

		kvArgs, err = client.getRandomCommitmentV2(shardID, common.PRVIDStr, len(inputCoins)*(privacy.RingSize-1))
		if err != nil {
			return nil, txHash, err
		}
		kvArgs[utils.MyIndices] = coinIndices

		txInitParam := tx_generic.NewTxPrivacyInitParams(&(senderWallet.KeySet.PrivateKey), paymentInfos, inputCoins, txFee, true, &common.PRVCoinID, nil, nil, kvArgs)
		tx := new(tx_ver2.Tx)
		err = tx.Init(txInitParam)
		if err != nil {
			return nil, "", fmt.Errorf("init txver2 error: %v", err)
		}

		txBytes, err := json.Marshal(tx)
		if err != nil {
			return nil, "", fmt.Errorf("cannot marshal txver2: %v", err)
		}

		base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)
		return []byte(base58CheckData), tx.Hash().String(), nil
	}
}

// SendRawTx sends submits a raw PRV transaction to the Incognito blockchain.
func (client *IncClient) SendRawTx(encodedTx []byte) error {
	responseInBytes, err := client.rpcServer.SendRawTx(string(encodedTx))
	if err != nil {
		return nil
	}

	err = rpchandler.ParseResponse(responseInBytes, nil)
	if err != nil {
		return err
	}

	return nil
}
