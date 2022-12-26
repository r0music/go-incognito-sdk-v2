package coin

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// OTAReceiver holds the data necessary to receive a coin with privacy.
// It is somewhat equivalent in usage with PaymentAddress.
type OTAReceiver struct {
	// PublicKey is the one-time public key of the receiving coin.
	PublicKey crypto.Point

	// TxRandom is for the receiver to recover the receiving information.
	TxRandom TxRandom

	// SharedSecrets are for the sender to mask the amount as well as the asset type of the sending coins.
	// SharedSecrets = []crypto.Point{sharedOTAPoint, sharedConcealPoint}:
	//	- sharedOTAPoint: used for generating the one-time address and concealing the assetID.
	//	- sharedConcealPoint: used for concealing the amount.
	// For non-privacy transactions, this field can be omitted.
	// SharedSecrets []crypto.Point `json:"SharedSecrets,omitempty"`
}

// IsValid checks the validity of this OTAReceiver (all referenced Points must be valid).
// Note that some sanity checks are already done when unmarshalling
func (receiver OTAReceiver) IsValid() bool {
	_, err := receiver.TxRandom.GetTxConcealRandomPoint()
	if err != nil {
		return false
	}
	_, err = receiver.TxRandom.GetTxOTARandomPoint()
	if err != nil {
		return false
	}
	// if len(receiver.SharedSecrets) > 0 {
	// 	if len(receiver.SharedSecrets) != 2 {
	// 		return false
	// 	}
	// 	if !receiver.SharedSecrets[0].PointValid() || !receiver.SharedSecrets[1].PointValid() {
	// 		return false
	// 	}
	// }
	return receiver.PublicKey.PointValid()
}

// IsConcealable checks if the OTAReceiver supports full privacy.
// func (receiver OTAReceiver) IsConcealable() bool {
// 	return len(receiver.SharedSecrets) == 2
// }

// GetPublicKey returns the base58-encoded PublicKey of an OTAReceiver.
func (receiver OTAReceiver) GetPublicKey() string {
	return base58.Base58Check{}.Encode(receiver.PublicKey.ToBytesS(), 0)
}

// GetTxRandom returns the base58-encoded GetTxRandom of an OTAReceiver.
func (receiver OTAReceiver) GetTxRandom() string {
	return base58.Base58Check{}.Encode(receiver.TxRandom.Bytes(), 0)
}

func (recv *OTAReceiver) From(addr key.PaymentAddress, targetSenderShardID, cptype int, withConceal bool) error {
	if recv == nil {
		return errors.New("OTAReceiver not initialized")
	}

	targetShardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])
	if targetSenderShardID == -1 {
		targetSenderShardID = int(targetShardID)
	}

	otaRand := crypto.RandomScalar()
	concealRand := (&crypto.Scalar{}).FromUint64(0)
	if withConceal {
		concealRand = crypto.RandomScalar()
	}

	index := uint32(0)
	publicOTA := addr.GetOTAPublicKey()
	if publicOTA == nil {
		return errors.New("Missing public OTA in payment address")
	}
	publicSpend := addr.GetPublicSpend()
	rK := (&crypto.Point{}).ScalarMult(publicOTA, otaRand)
	for i := MaxTriesOTA; i > 0; i-- {
		index++
		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := (&crypto.Point{}).ScalarMultBase(hash)
		publicKey := (&crypto.Point{}).Add(HrKG, publicSpend)

		pkb := publicKey.ToBytesS()
		senderShardID, recvShardID, coinPrivacyType, _ := DeriveShardInfoFromCoin(pkb)
		if recvShardID == int(targetShardID) && senderShardID == targetSenderShardID && coinPrivacyType == cptype {
			otaRandomPoint := (&crypto.Point{}).ScalarMultBase(otaRand)
			concealRandomPoint := (&crypto.Point{}).ScalarMultBase(concealRand)
			recv.PublicKey = *publicKey
			recv.TxRandom = *NewTxRandom()
			recv.TxRandom.SetTxOTARandomPoint(otaRandomPoint)
			recv.TxRandom.SetTxConcealRandomPoint(concealRandomPoint)
			recv.TxRandom.SetIndex(index)
			return nil
		}
	}
	return fmt.Errorf("Cannot generate OTAReceiver after %d attempts", MaxTriesOTA)
}

func (receiver *OTAReceiver) FromAddress(addr key.PaymentAddress) error {
	if receiver == nil {
		return errors.New("OTAReceiver not initialized")
	}

	targetShardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])
	otaRand := crypto.RandomScalar()
	concealRand := crypto.RandomScalar()

	index := uint32(0)
	publicOTA := addr.GetOTAPublicKey()
	if publicOTA == nil {
		return errors.New("Missing public OTA in payment address")
	}
	publicSpend := addr.GetPublicSpend()
	rK := (&crypto.Point{}).ScalarMult(publicOTA, otaRand)
	for i := MaxTriesOTA; i > 0; i-- {
		index++
		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := (&crypto.Point{}).ScalarMultBase(hash)
		publicKey := (&crypto.Point{}).Add(HrKG, publicSpend)

		pkb := publicKey.ToBytesS()
		senderShardID, recvShardID, coinPrivacyType, _ := DeriveShardInfoFromCoin(pkb)
		if recvShardID == int(targetShardID) && senderShardID == int(targetShardID) && coinPrivacyType == PrivacyTypeMint {
			otaRandomPoint := (&crypto.Point{}).ScalarMultBase(otaRand)
			concealRandomPoint := (&crypto.Point{}).ScalarMultBase(concealRand)
			receiver.PublicKey = *publicKey
			receiver.TxRandom = *NewTxRandom()
			receiver.TxRandom.SetTxOTARandomPoint(otaRandomPoint)
			receiver.TxRandom.SetTxConcealRandomPoint(concealRandomPoint)
			receiver.TxRandom.SetIndex(index)
			return nil
		}
	}
	return fmt.Errorf("Cannot generate OTAReceiver after %d attempts", MaxTriesOTA)
}

// FromString returns a new OTAReceiver parsed from the input string,
// or error on failure
func (receiver *OTAReceiver) FromString(data string) error {
	raw, _, err := base58.Base58Check{}.Decode(data)
	if err != nil {
		return err
	}
	err = receiver.SetBytes(raw)
	if err != nil {
		return err
	}
	return nil
}

// String marshals the OTAReceiver, then encodes it with base58.
// By default, an OTAReceiver will only support receiving assets in a non-private transaction. Set `isConcealable = true`
// to enable receiving assets in a private transaction.
func (receiver OTAReceiver) String() string {
	return base58.Base58Check{}.NewEncode(receiver.Bytes(), common.ZeroByte)
}

// Bytes returns a byte-encoded form of an OTAReceiver.
// By default, an OTAReceiver will only support receiving assets in a non-private transaction.
func (receiver OTAReceiver) Bytes() []byte {
	rawBytes := []byte{wallet.PrivateReceivingAddressType}
	rawBytes = append(rawBytes, receiver.PublicKey.ToBytesS()...)
	rawBytes = append(rawBytes, receiver.TxRandom.Bytes()...)
	return rawBytes
}

func (receiver *OTAReceiver) SetBytes(b []byte) error {
	if len(b) == 0 {
		return errors.New("Not enough bytes to parse ReceivingAddress")
	}
	if receiver == nil {
		return errors.New("OTAReceiver not initialized")
	}
	keyType := b[0]
	switch keyType {
	case wallet.PrivateReceivingAddressType:
		buf := make([]byte, 32)
		copy(buf, b[1:33])
		pk, err := (&crypto.Point{}).FromBytesS(buf)
		if err != nil {
			return err
		}
		receiver.PublicKey = *pk
		txr := NewTxRandom()
		// SetBytes() will perform length check
		err = txr.SetBytes(b[33:])
		if err != nil {
			return err
		}
		receiver.TxRandom = *txr
		return nil
	default:
		return errors.New("Unrecognized prefix for ReceivingAddress")
	}
}

// MarshalJSON returns a non-private byte-sequence representation of an OTAReceiver.
func (receiver OTAReceiver) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiver.String())
}

func (receiver *OTAReceiver) UnmarshalJSON(raw []byte) error {
	var encodedString string
	err := json.Unmarshal(raw, &encodedString)
	if err != nil {
		return err
	}
	var temp OTAReceiver
	err = temp.FromString(encodedString)
	if err != nil {
		return err
	}
	*receiver = temp
	return nil
}

// GetShardIDs returns a pair of (sendingShard, receivingShard) of an OTAReceiver.
func (receiver OTAReceiver) GetShardIDs() (byte, byte) {
	pkb := receiver.PublicKey.ToBytesS()
	return common.GetShardIDsFromPublicKey(pkb)
}

// FromCoinParams generates an OTAReceiver from the given CoinParams.
func (receiver *OTAReceiver) FromCoinParams(p *CoinParams) error {
	if receiver == nil {
		return fmt.Errorf("OTAReceiver not initialized")
	}

	addr := p.PaymentInfo.PaymentAddress

	receiverShardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])
	otaRand := crypto.RandomScalar()
	concealRand := crypto.RandomScalar()

	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := addr.GetOTAPublicKey()
	if publicOTA == nil {
		return fmt.Errorf("missing public OTA in payment address")
	}
	publicSpend := addr.GetPublicSpend()
	rK := (&crypto.Point{}).ScalarMult(publicOTA, otaRand)
	for i := MaxTriesOTA; i > 0; i-- {
		index++
		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := (&crypto.Point{}).ScalarMultBase(hash)
		publicKey := (&crypto.Point{}).Add(HrKG, publicSpend)

		tmpSenderShardID, tmpReceiverShardID, tmpCoinType, _ := DeriveShardInfoFromCoin(publicKey.ToBytesS())
		if tmpReceiverShardID == int(receiverShardID) && tmpSenderShardID == p.SenderShardID && tmpCoinType == p.CoinPrivacyType {
			otaRandomPoint := (&crypto.Point{}).ScalarMultBase(otaRand)
			concealRandomPoint := (&crypto.Point{}).ScalarMultBase(concealRand)
			receiver.PublicKey = *publicKey
			receiver.TxRandom = *NewTxRandom()
			receiver.TxRandom.SetTxOTARandomPoint(otaRandomPoint)
			receiver.TxRandom.SetTxConcealRandomPoint(concealRandomPoint)
			receiver.TxRandom.SetIndex(index)
			return nil
		}
	}
	return fmt.Errorf("cannot generate OTAReceiver after %d attempts", MaxTriesOTA)
}
func (receiver OTAReceiver) DeriveShardID() (int, int, int, error) {
	return DeriveShardInfoFromCoin(receiver.PublicKey.ToBytesS())
}
