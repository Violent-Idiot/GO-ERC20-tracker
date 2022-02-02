package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
}

type Transfer struct {
	From   string
	To     string
	Amount string
}

var TransferArray []Transfer

// type LogApproval struct {
// 	Owner   common.Address
// 	Spender common.Address
// 	Amount  *big.Int
// }

func main() {
	endpoint := os.Getenv("INFURA")
	client, err := ethclient.DialContext(context.Background(), endpoint)
	if err != nil {
		log.Fatal(err)
	}
	addr := common.HexToAddress("0x6f40d4a6237c257fff2db00fa0510deeecd303eb")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(14095312), //14095312
		ToBlock:   big.NewInt(14119242),
		Addresses: []common.Address{addr},
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("./abi")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	parsed, err := abi.JSON(file)
	if err != nil {
		log.Fatal(err)
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	// LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	// logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)
	for _, vLog := range logs {
		// fmt.Printf("\nLog Block Number: %d\n", vLog.BlockNumber)
		// fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			// fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			err := parsed.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			temp := Transfer{
				From:   transferEvent.From.Hex(),
				To:     transferEvent.To.Hex(),
				Amount: transferEvent.Amount.String(),
			}
			TransferArray = append(TransferArray, temp)
			// fmt.Printf("From: %s\n", transferEvent.From.Hex())
			// fmt.Printf("To: %s\n", transferEvent.To.Hex())
			// fmt.Printf("Tokens: %s\n", transferEvent.Amount.String())
			// case logApprovalSigHash.Hex():
			// 	fmt.Printf("Log Name: Approval\n")

			// 	var approvalEvent LogApproval

			// 	err := parsed.UnpackIntoInterface(&approvalEvent, "Approval", vLog.Data)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}

			// 	approvalEvent.Owner = common.HexToAddress(vLog.Topics[1].Hex())
			// 	approvalEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())

			// 	fmt.Printf("Token Owner: %s\n", approvalEvent.Owner.Hex())
			// 	fmt.Printf("Spender: %s\n", approvalEvent.Spender.Hex())
			// 	fmt.Printf("Tokens: %s\n\n", approvalEvent.Amount.String())
		}
	}
	sort.Slice(TransferArray, func(i, j int) bool {
		return TransferArray[i].Amount > TransferArray[j].Amount
	})
	TransferArray = TransferArray[:15]
	for _, item := range TransferArray {
		fmt.Printf("\nFrom:- %s\n", item.From)
		fmt.Printf("To:- %s\n", item.To)
		fmt.Printf("Amount:- %s\n\n", item.Amount)
	}
	defer client.Close()
}
