package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"sort"
	"sync"
	"time"

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
	start := time.Now()
	defer func() {
		fmt.Println("Exection time: ", time.Since(start))
	}()
	os.Setenv("INFURA", "wss://mainnet.infura.io/ws/v3/eb979022577d4b55b620e583cc58ba72")
	choose := 0
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
	limit := 10000
	tempUpperBound := lowerBound + limit
	var logs []types.Log
	flag := false
	paginate := (upperBound - lowerBound) / limit
	fmt.Println(paginate)
	wg := sync.WaitGroup{}
	if choose == 0 {

		for i := 0; i < paginate; i++ {

			fmt.Println(i)
			wg.Add(1)
			lower := lowerBound + (i * limit)
			tempUpper := tempUpperBound + (i * limit)
			// fmt.Println(lower, tempUpper)
			if tempUpper > upperBound {
				tempUpper = upperBound
				flag = true
				// break
			}
			if flag {
				break
			}
			go func(lower, tempUpper int) {
				query := ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(lower)), //14095312
					ToBlock:   big.NewInt(int64(tempUpper)),
					Addresses: []common.Address{addr},
					Topics:    [][]common.Hash{{logTransferSigHash}},
				}

				templogs, err := client.FilterLogs(context.Background(), query)
				if err != nil {
					log.Fatal(err)
				}
				logs = append(logs, templogs...)
				// lowerBound += limit
				// tempUpperBound += limit
				fmt.Println(lower, tempUpper)
				// if tempUpperBound >= upperBound {
				// 	tempUpperBound = upperBound
				// 	flag = true
				// }
				// if flag {
				// 	break
				// }
				wg.Done()
			}(lower, tempUpper)

		}
		wg.Wait()
	}
	if choose == 1 {

		for {

			query := ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(lowerBound)), //14095312
				ToBlock:   big.NewInt(int64(tempUpperBound)),
				Addresses: []common.Address{addr},
				Topics:    [][]common.Hash{{logTransferSigHash}},
			}

			templogs, err := client.FilterLogs(context.Background(), query)
			if err != nil {
				log.Fatal(err)
			}
			logs = append(logs, templogs...)
			lowerBound += limit
			tempUpperBound += limit
			fmt.Println(lowerBound, tempUpperBound)
			if tempUpperBound >= upperBound {
				tempUpperBound = upperBound
				flag = true
			}
			if flag {
				break
			}
		}
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

	H := make(map[string]float64)
	log.Println(len(logs))
	init := true

	for index, vLog := range logs {
		fmt.Println("mapping", index)
		if vLog.Topics[0].Hex() == logTransferSigHash.Hex() {

			var transferEvent LogTransfer

			err := parsed.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			// fmt.Println(transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Amount.String())
			bFloat, _ := new(big.Float).SetString(transferEvent.Amount.String())
			floatValue := new(big.Float).Quo(bFloat, big.NewFloat(math.Pow10(18)))
			fValue, _ := floatValue.Float64()
			from := H[transferEvent.From.Hex()]
			to := H[transferEvent.To.Hex()]

			if init {
				from = 0
				to += fValue
				fmt.Println(from, to)
				init = false

			} else {
				from -= fValue
				to += fValue
			}
			// fmt.Println(from, to, fValue)
			H[transferEvent.From.Hex()] = from
			H[transferEvent.To.Hex()] = to
		}
	}
	type kv struct {
		Key   string
		Value float64
	}

	var ss []kv

	for k, v := range H {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	ss = ss[:15]
	fmt.Println()
	fmt.Println("Top 15 INST Holders:-")
	fmt.Println()
	for _, kv := range ss {
		fmt.Printf("%s %f\n", kv.Key, kv.Value)
	}

	defer client.Close()
}
