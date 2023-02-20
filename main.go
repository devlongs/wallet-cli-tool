package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Parse command line flags
	destPtr := flag.String("dest", "", "destination address to transfer ether")
	amountPtr := flag.String("amount", "", "amount of ether to transfer (in wei)")
	flag.Parse()

	// Connect to an Ethereum client
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/your-project-id")
	if err != nil {
		panic(err)
	}

	// Generate a new Ethereum account
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("failed to convert public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Print the account address
	fmt.Println("New account address:", address.Hex())

	// Send a transaction from the new account to the destination address
	toAddress := common.HexToAddress(*destPtr)
	amount, ok := new(big.Int).SetString(*amountPtr, 10)
	if !ok {
		panic("failed to parse amount")
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		panic(err)
	}
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		panic(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sent %v ETH from %v to %v\n", amount, address.Hex(), toAddress.Hex())
}
