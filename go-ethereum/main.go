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
	value := big.NewInt(1000000000000000000)
	err = token.TransferETH(privateKeyHardhat, addressHardhat, privateKey, address, value, client)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy the smart contract
	tokenAddress, instance := token.Deploy(address, privateKey, client)

	// Check information are correct
	if err != nil {
		log.Fatal(err)
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	supply, err := instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	balanceA, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Token's name: %s\n", name)
	fmt.Printf("Token's symbol: %s\n", symbol)
	fmt.Printf("Token's supply: %s\n", supply)

	fmt.Printf("A's Balance is: %s\n", new(big.Int).Div(balanceA, big.NewInt(1000000000000000000)))

	// Create another address
	// Generate new account B
	privateKeyTo, addressTo := accounts.New()

	// Transfer MTK token
	value = big.NewInt(0)
	err = token.Transfer(privateKey, address, privateKeyTo, addressTo, tokenAddress, value, instance, client)
	if err != nil {
		log.Fatal(err)
	}

	balanceAAfter, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("A's Balance after the transaction is: %s\n", new(big.Int).Div(balanceAAfter, big.NewInt(1000000000000000000)))
}
