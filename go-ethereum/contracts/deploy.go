package token

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Deploy(address common.Address, privateKey *ecdsa.PrivateKey, client *ethclient.Client) (tokenAdress common.Address, instance *Token) {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction with the private key of the sender
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	tokenAddress, tx, instance, err := DeployToken(auth, client)
	if err != nil {
		log.Fatal(err)
	}
	_ = tx
	fmt.Printf("Token address is: %s\n", tokenAddress.Hex())
	return tokenAddress, instance
}
