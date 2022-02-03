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
	"github.com/ethereum/go-ethereum/core/types"
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

func main() {
	os.Setenv("INFURA", "wss://mainnet.infura.io/ws/v3/eb979022577d4b55b620e583cc58ba72")
	endpoint := os.Getenv("INFURA")
	client, err := ethclient.DialContext(context.Background(), endpoint)
	if err != nil {
		log.Fatal(err)
	}
	addr := common.HexToAddress("0x6f40d4a6237c257fff2db00fa0510deeecd303eb")
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	upperBound := 14119242
	// lowerBound := 14101404
	lowerBound := 12183236
	// TotalBlocks := upperBound - lowerBound
	// tempUpperBound := upperBound
	limit := 10000
	tempUpperBound := lowerBound + limit
	// div := 0
	var logs []types.Log
	flag := false
	// MidBlock := 0
	for {

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(lowerBound)), //14095312
			ToBlock:   big.NewInt(int64(tempUpperBound)),
			Addresses: []common.Address{addr},
			Topics:    [][]common.Hash{{logTransferSigHash}},
		}

		templogs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			// if err != nil {
			// TotalBlocks /= 2
			// fmt.Println(TotalBlocks)
			// tempUpperBound = lowerBound + TotalBlocks
			// fmt.Println(lowerBound, tempUpperBound)
			// continue
			log.Fatal(err)
			// div += 1
			// MidBlock = TotalBlocks
		}
		// else {
		// fmt.Println("Here")
		logs = append(logs, templogs...)
		// log.Println("loggin")
		lowerBound += limit
		// lowerBound += TotalBlocks
		tempUpperBound += limit
		// tempUpperBound += TotalBlocks
		fmt.Println(lowerBound, tempUpperBound)
		if tempUpperBound >= upperBound {
			tempUpperBound = upperBound
			flag = true
		}
		if flag {
			break
		}
		// }
	}
	fmt.Println("here")
	file, err := os.Open("./abi")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	parsed, err := abi.JSON(file)
	if err != nil {
		log.Fatal(err)
	}

	H := make(map[string]*big.Int)
	log.Println(len(logs))
	// ctx := context.Background()
	init := true

	for index, vLog := range logs {
		// fmt.Printf("\nLog Block Number: %d\n", vLog.BlockNumber)
		// fmt.Printf("Log Index: %d\n", vLog.Index)

		fmt.Println("mapping", index)
		// log.Print(vLog.Topics[0].Hex(), logTransferSigHash.Hex())
		if vLog.Topics[0].Hex() == logTransferSigHash.Hex() {

			var transferEvent LogTransfer

			err := parsed.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			// account := common.HexToAddress(transferEvent.From.Hex())
			fmt.Println(transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Amount.String())
			bigInt := new(big.Int)
			bInt, _ := bigInt.SetString(transferEvent.Amount.String(), 10)
			// break
			// bInt
			// fmt.Println(bInt)
			// break
			// fmt.Println(H[transferEvent.From.Hex()])
			// break
			// if big.NewInt(0).Sub(H[transferEvent.From.Hex()], bInt) < big.NewInt(0){
			// if H[transferEvent.From.Hex()].Cmp(bInt) == -1 {
			from := H[transferEvent.From.Hex()]
			to := H[transferEvent.To.Hex()]
			if to == nil {
				to = big.NewInt(0)
			}
			if from == nil {
				from = big.NewInt(0)
			}
			// to := big.NewInt(0)
			fmt.Println(init)

			if init {
				// yoamn = big.NewInt(0)
				fmt.Println("here")
				// fmt.Println(to)
				from = big.NewInt(0)
				to = to.Add(to, bInt)
				// H[transferEvent.From.Hex()] = from

				// H[transferEvent.To.Hex()] = to

				fmt.Println(from, to)

				init = false

			} else {
				// H[transferEvent.From.Hex()] = big.NewInt(0).Sub(H[transferEvent.From.Hex()], bInt)
				// H[transferEvent.To.Hex()] = big.NewInt(0).Add(H[transferEvent.To.Hex()], bInt)

				from = from.Sub(from, bInt)
				to = to.Add(to, bInt)
				fmt.Println(from, to, bInt)

			}
			H[transferEvent.From.Hex()] = from

			H[transferEvent.To.Hex()] = to
			// fmt.Println(H[transferEvent.To.Hex()])
			// break
			// bal, _ := client.BalanceAt(ctx, account, nil)
			// fmt.Println(transferEvent.From.Hex())
			// H[transferEvent.From.Hex()] = int(bal.Uint64())
			// fmt.Println(transferEvent.From.Hex(), bal)
			// temp := Transfer{
			// 	From:   transferEvent.From.Hex(),
			// 	To:     transferEvent.To.Hex(),
			// 	Amount: transferEvent.Amount.String(),
			// }
			// TransferArray = append(TransferArray, temp)
		}
	}
	type kv struct {
		Key   string
		Value *big.Int
	}

	var ss []kv

	for k, v := range H {
		ss = append(ss, kv{k, v})
	}

	sort.SliceStable(ss, func(i, j int) bool {
		// return ss[i].Value > ss[j].Value
		// fmt.Println(ss[i].Value, ss[j].Value, ss[i].Value.Cmp(ss[j].Value) > 1)
		return ss[i].Value.Cmp(ss[j].Value) > 1
	})
	ss = ss[:15]

	for _, kv := range ss {
		// temp := new(big.Float)
		// temp.SetUint64(kv.Value)
		// value := new(big.Float).Quo(temp, big.NewFloat(math.Pow10(18)))
		fmt.Printf("%s %d\n", kv.Key, kv.Value)
	}

	defer client.Close()
}
