package main

import (
	"bytes"
	"fmt"
	goSdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	// "github.com/ontio/ontology/core/types"
	"os"
	"strconv"
)

const ROUTING_NUM = 10

func main() {
	defer os.RemoveAll("./Log")
	log.InitLog(log.InfoLog)

	ontSdk := goSdk.NewOntologySdk()
	// get the FROM wallet
	var wallet, _ = ontSdk.OpenWallet("wallet-admin.dat")
	// get the FROM account inside the wallet
	admin, err := wallet.GetAccountByAddress("AT6HVkWJ3HdAKidZEH1kbYZ7y896mASQU2", []byte("111"))
	if err != nil {
		fmt.Println("get admin err:", err)
	}

	// get the TO account wallet
	wallet, _ = ontSdk.OpenWallet("wallet-account.dat")
	toAcc, err := wallet.GetAccountByAddress("AaZjCJnyLTFt4rCvBw7kagZkE1CLzLQUHy", []byte("111"))
	if err != nil {
		fmt.Println("get account err:", err)
	}

	txNum, _ := strconv.Atoi(os.Args[1])
	txNum = txNum * 100000
	if txNum > 2<<32 {
		txNum = 2 << 32
	}
	exitChan := make(chan int)
	txNumPerRoutine := txNum / ROUTING_NUM
	for i := 0; i < ROUTING_NUM; i++ {
		go func(nonce uint32, routineIndex int) {
			for j := 0; j < txNumPerRoutine; j++ {
				txHash, txContent := genTransfer(ontSdk, admin, toAcc, 1, nonce)
				nonce++
				fmt.Print(txHash, ",", txContent, "\n")
			}
			exitChan <- 1
		}(uint32(txNumPerRoutine*i), i)
	}
	for i := 0; i < ROUTING_NUM; i++ {
		<-exitChan
	}
}

func genTransfer(ontSdk *goSdk.OntologySdk, from *goSdk.Account, to *goSdk.Account, value uint64, nonce uint32) (string, string) {
	mtx, err := ontSdk.Native.Ont.NewTransferTransaction(0, 100000, from.Address, to.Address, value) //func (this *Ont) NewTransferTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error)
	if err != nil {
		return "", ""
	}
	mtx.Nonce = nonce
	
	err = ontSdk.SignToTransaction(mtx, from)
	if err != nil {
		return "", ""
	}

	txbf := new(bytes.Buffer)
	//zcs := common.NewZeroCopySink(nil)
	tx, err := mtx.IntoImmutable()
	if err != nil {
		fmt.Println("Convert to immutable error.")
		os.Exit(1)
	}
	if err := tx.Serialize(txbf); err != nil {
		fmt.Println("Serialize transaction error.")
		os.Exit(1)
	}
	hash := tx.Hash()
	return hash.ToHexString(), common.ToHexString(txbf.Bytes())
}