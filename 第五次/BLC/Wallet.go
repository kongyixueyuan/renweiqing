package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)
const version = byte(0x00)
const addressChecksumLen = 4

// 钱包
type Wallet struct {
	// 私钥
	Private ecdsa.PrivateKey
	// 公钥
	PublicKey []byte
}

func IsValidForAddress(address []byte) bool  {
	// 将地址 解码 25位字节数组
	version_public_checksumBytes := Base58Decode(address)
    // 获取后4位
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressChecksumLen : ]
    // 获取前21位
	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes) - addressChecksumLen]
	// sha256(sha256(PubKeyHash))[:addressChecksumLen]
	checkBytes := CheckSum(version_ripemd160)
	// 返回比对的结果
	return bytes.Compare(checkBytes,checkSumBytes) == 0
}
// Base58Encode(version + ripemd160(sha256(PubKey)) + sha256(sha256(PubKeyHash))[:addressChecksumLen])
func (w *Wallet) GetAddress() []byte  {
	// 1. hash160
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	// 将version和ripemd160Hash拼接到一起
	version_ripemd160Hash := append([]byte{version},ripemd160Hash...)
	// 获取验证值
	checkSumBytes := CheckSum(version_ripemd160Hash)
	// 将 version、ripemd160Hash、checkSum拼接到一起
	bytes := append(version_ripemd160Hash,checkSumBytes...)
	// 返回base58加密串
	return Base58Encode(bytes)
}
// sha256(sha256(PubKeyHash))[:addressChecksumLen]
func CheckSum(b []byte) []byte  {
	hash1 := sha256.Sum256(b)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}
// ripemd160(sha256(PubKey))
func Ripemd160Hash(publicKey []byte) []byte  {
	// 1. 256
	//hash256 := sha256.New()
	//hash256.Write(publicKey)
	//hash := hash256.Sum(nil)
	hash := sha256.Sum256(publicKey)

	// 2. 160
	r160 := ripemd160.New()
	r160.Write(hash[:])
	return r160.Sum(nil)
}

// 创建钱包
func NewWallet() *Wallet {

	privateKey,publicKey := newKeyPair()

	return &Wallet{privateKey,publicKey}
}


// 通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey,[]byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}