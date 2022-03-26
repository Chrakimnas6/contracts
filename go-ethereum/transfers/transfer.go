package transfers

import (
	token "chrakimnas6/go-ethereum/contracts"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func TransferETH(privateKeyFrom *ecdsa.PrivateKey, addressFrom common.Address,
	privateKeyTo *ecdsa.PrivateKey, addressTo common.Address,
	value *big.Int, client *ethclient.Client) (err error) {

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

func TransferToken(privateKeyFrom *ecdsa.PrivateKey, addressFrom common.Address,
	privateKeyTo *ecdsa.PrivateKey, addressTo common.Address,
	tokenAddress common.Address, value *big.Int, instance *token.Token, client *ethclient.Client) (err error) {

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(addressTo.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &addressTo,
		Data: data,
	})
	if err != nil {
		return err
	}

	fmt.Println(gasLimit)

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

	tx := types.NewTransaction(nonce, tokenAddress, value, 300000, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyFrom)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	balanceNew, err := instance.BalanceOf(&bind.CallOpts{}, addressTo)
	if err != nil {
		return err
	}
	fmt.Printf("Balance is %s\n", balanceNew)

	return nil
}
