package BLC

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"bytes"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Rwq_Wallet struct {
	Rwq_PrivateKey ecdsa.PrivateKey
	Rwq_PublicKey  []byte
}

// 重成新的钱包
func Rwq_NewWallet() *Rwq_Wallet {
	private, public := rwq_newKeyPair()
	wallet := Rwq_Wallet{private, public}
	return &wallet
}

// 获取钱包地址
func (w Rwq_Wallet) Rwq_GetAddress() []byte {
	pubKeyHash := Rwq_HashPubKey(w.Rwq_PublicKey)
	versionPayload := append([]byte{version}, pubKeyHash...)
	checksum := rwq_checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	return Base58Encode(fullPayload)
}

// 将公钥先进行Hash256 再进行 RIPEMD160 Hash
func Rwq_HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	return RIPEMD160Hasher.Sum(nil)
}

// 验证钱包地址是否有效
// 地址解码后，将前21个字符两次hash后
// 取前四位跟 解码后的地址后四位对比
func Rwq_ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := rwq_checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// 两次hash256获取校验值，hash的前4个
func rwq_checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

// 生成新的私钥和公钥
func rwq_newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
