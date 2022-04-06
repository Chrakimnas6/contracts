package main

import (
	"chrakimnas6/go-ethereum/accounts"
	token "chrakimnas6/go-ethereum/contracts"
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	privateKeyHardhat, addressHardhat := accounts.GetHardhatAddress()

	// Check ETH balance of the address
	balance, err := client.BalanceAt(context.Background(), addressHardhat, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ETH balance of the Hardhat's local node: %s\n", balance)

	// Generate new account A
	privateKey, address := accounts.New()

	// Transferring ETH to A's address
	value := big.NewInt(1000000000000000000) // 1 eth
	err = token.TransferETH(privateKeyHardhat, addressHardhat, address, value, client)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy the smart contract
	tokenAddress, instance := token.Deploy(address, privateKey, client)
	_ = tokenAddress

	// Check information are correct
	token.CheckInformation(instance, address)

	// Mint extra 10000 tokens to A's address so he will have 20000 tokens in total
	value = new(big.Int)
	value.SetString("10000000000000000000000", 10) // 10000 tokens
	err = token.MintUsingAPI(privateKey, address, value, instance, client)
	if err != nil {
		log.Fatal(err)
	}

	// Check information are correct
	token.CheckInformation(instance, address)

	// // Check information are correct
	// token.CheckInformation(instance, address)

	// Create another address
	// Generate new account B
	privateKeyTo, addressTo := accounts.New()
	_ = privateKeyTo

	// Transfer MTK token
	value = new(big.Int)
	value.SetString("1000000000000000000000", 10) // 1000 tokens

	// Transfering using API
	// err = token.TransferUsingAPI(privateKey, address, addressTo, value, instance, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Normal transfer
	err = token.Transfer(privateKey, address, addressTo, tokenAddress, value, instance, client)
	if err != nil {
		log.Fatal(err)
	}

	balanceAAfter, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("A's Balance after the transaction is: %s\n", new(big.Int).Div(balanceAAfter, big.NewInt(1000000000000000000)))
}
