package accounts

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Get account from Hardhat local node
func GetHardhatAddress() (privateKey *ecdsa.PrivateKey, address common.Address) {
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}
	address = deriveAddress(privateKey)
	return privateKey, address
}

// Generate new account
func GenerateNewAccount() (privateKey *ecdsa.PrivateKey, address common.Address) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	address = deriveAddress(privateKey)
	return privateKey, address
}

// Derive address from private key
func deriveAddress(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address
}
