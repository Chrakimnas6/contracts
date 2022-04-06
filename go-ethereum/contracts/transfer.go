package token

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func MintUsingAPI(privateKey *ecdsa.PrivateKey, address common.Address,
	value *big.Int, instance *Token, client *ethclient.Client) (err error) {

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	tx, err := instance.Mint(auth, address, value)
	if err != nil {
		return err
	}
	_ = tx
	return nil
}

func TransferETH(privateKeyFrom *ecdsa.PrivateKey, addressFrom common.Address,
	addressTo common.Address, value *big.Int, client *ethclient.Client) (err error) {

	nonce, err := client.PendingNonceAt(context.Background(), addressFrom)
	if err != nil {
		return err
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, addressTo, value, gasLimit, gasPrice, nil)

	// Sign the transaction with the private key of the sender
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	//
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyFrom)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	return nil
}

func Transfer(privateKey *ecdsa.PrivateKey, addressFrom common.Address, addressTo common.Address,
	tokenAddress common.Address, value *big.Int, instance *Token, client *ethclient.Client) (err error) {

	// ERC-20 specification
	transferFnSignature := []byte("transfer(address,uint256)")
	// Generate the Keccak256 hash of the function signature
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	// Left pad 32 bytes the address we'are sending tokens to
	paddedAddress := common.LeftPadBytes(addressTo.Bytes(), 32)

	// Also left padding 32 bits for the amount
	paddedValue := common.LeftPadBytes(value.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedValue...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &addressTo,
		Data: data,
	})
	if err != nil {
		return err
	}
	// Not enough
	_ = gasLimit

	nonce, err := client.PendingNonceAt(context.Background(), addressFrom)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if nil != err {
		return err
	}

	chainID, err := client.NetworkID(context.Background())
	if nil != err {
		return err
	}
	// Gas limit is not enough with the estimated one, so set to 300000 here for now
	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), 300000, gasPrice, data)

	// Sign the transaction with the private key of the sender
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	balanceB, err := instance.BalanceOf(&bind.CallOpts{}, addressTo)
	if err != nil {
		return err
	}

	fmt.Printf("B's Balance is %s\n", new(big.Int).Div(balanceB, big.NewInt(1000000000000000000)))

	return nil
}

func TransferUsingAPI(privateKey *ecdsa.PrivateKey, addressFrom common.Address, addressTo common.Address,
	value *big.Int, instance *Token, client *ethclient.Client) (err error) {
	nonce, err := client.PendingNonceAt(context.Background(), addressFrom)
	if err != nil {
		return err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	tx, err := instance.Transfer(auth, addressTo, value)
	if err != nil {
		return err
	}
	_ = tx
	return nil
}
