package main

import (
	"chrakimnas6/go-ethereum/accounts"
	token "chrakimnas6/go-ethereum/contracts"
	"chrakimnas6/go-ethereum/transfers"
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
	fmt.Printf("ETH balance of the local node: %s\n", balance)

	// Generate new account A
	privateKey, address := accounts.GenerateNewAccount()

	// Transferring ETH to the address
	value := big.NewInt(1000000000000000000)
	err = transfers.TransferETH(privateKeyHardhat, addressHardhat, privateKey, address, value, client)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy the smart contract
	tokenAddress, instance := token.DeploySmartContract(address, privateKey, client)

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

	fmt.Printf("name: %s\n", name)
	fmt.Printf("symbol: %s\n", symbol)
	fmt.Printf("supply: %s\n", supply)

	// Create another address
	// Generate new account B
	privateKeyTo, addressTo := accounts.GenerateNewAccount()
	_ = privateKeyTo

	// Transfer MTK token
	value = big.NewInt(0)
	err = transfers.TransferToken(privateKey, address, privateKeyTo, addressTo, tokenAddress, value, instance, client)
	if err != nil {
		log.Fatal(err)
	}
}
