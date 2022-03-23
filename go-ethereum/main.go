package main

import (
	token "chrakimnast6/go-ethereum/contracts"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to Hardhat RPC host
	client, err := ethclient.Dial("http://127.0.0.1:8545/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")

	// Address from Hardhat's local node
	privateKeyHardhat, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}
	publicKeyHardhat := privateKeyHardhat.Public()
	publicKeyECDSAHardHat, ok := publicKeyHardhat.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	addressHardhat := crypto.PubkeyToAddress(*publicKeyECDSAHardHat)
	fmt.Println(addressHardhat)
	balance, err := client.BalanceAt(context.Background(), addressHardhat, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	// Generate new wallet
	// Generate random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Convert it to bytes
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Convert it to hexadecimal string, strip off 0x
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

	// public key is derived from private key
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// strip off the 0x and the first 2 characters 04 which is always the EC prefix
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println(address)

	// Transferring ETH to the address
	nonce, err := client.PendingNonceAt(context.Background(), addressHardhat)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, address, value, gasLimit, gasPrice, nil)

	// Sign the transaction with the private key of the sender
	chainID, err := client.NetworkID(context.Background())
	if nil != err {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyHardhat)
	if err != nil {
		log.Fatal(nil)
	}

	// Broadcast the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy the smart contract
	nonce, err = client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := token.DeployToken(auth, client)
	_ = tx

	if err != nil {
		log.Fatal(err)
	}

	_ = instance

	// Mint token to the address
	// Loading the contract
	// tokenAddress := common.HexToAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3")
	// instance, err := token.NewToken(tokenAddress, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// name, err := instance.Name(&bind.CallOpts{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// symbol, err := instance.Symbol(&bind.CallOpts{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("name: %s\n", name)
	// fmt.Printf("symbol: %s\n", symbol)

	// nonce, err = client.PendingNonceAt(context.Background(), address)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(0)
	// auth.GasLimit = uint64(300000)
	// auth.GasPrice = gasPrice

	// _, err = instance.Mint(auth, address, big.NewInt(1000000000000000000))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = instance.Mint(auth, address, amount)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Keystore
	// var account accounts.Account
	// createKs(&account)
	// //importKs()
	// fmt.Println(account.Address)
	// balance, err := client.BalanceAt(context.Background(), account.Address, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(balance)

	// amount := new(big.Int)
	// amount.SetString("1000000000000000000000", 10)

	// _, err = instance.Mint(&bind.TransactOpts{}, account.Address, amount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// balance, err = client.BalanceAt(context.Background(), account.Address, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(balance)

}
