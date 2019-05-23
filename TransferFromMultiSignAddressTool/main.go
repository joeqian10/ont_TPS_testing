package main

import (
    "fmt"
    "github.com/ontio/ontology-crypto/keypair"
    ont "github.com/ontio/ontology-go-sdk"
    "github.com/ontio/ontology/common"
    "github.com/ontio/ontology/core/types"
    ont2 "github.com/ontio/ontology/smartcontract/service/native/ont"
    "strconv"
    "time"
)

var (
    OntSdk *ont.OntologySdk
    Passwd   = []byte("111")
    GasPrice = uint64(0)
    GasLimit = uint64(20000)
    //Addr1   *ont.Account
    addr1,_= common.AddressFromBase58("AT6HVkWJ3HdAKidZEH1kbYZ7y896mASQU2")
    //addr2,_= common.AddressFromBase58("AeM3iaPuVBn65sVf4HEAGCVE8HJP5Qshq4")
    //addr3,_= common.AddressFromBase58("ANys165m1NCPTDrPNWg5kHgiCsJSL4YgWP")
    //addr4,_= common.AddressFromBase58("ASAkfLiA38Tv3RnKg6p5ic4uyGQJMoXMX4")

    //addr5,_= common.AddressFromBase58("AKDjhL8LU5sjd3jFQbNQSkXCz9f2cavSog")
    //addr6,_= common.AddressFromBase58("AL5GixTHMGdn4iS2FMsULYq3R5c7WZy617")
    //addr7,_= common.AddressFromBase58("AXExefyGAFrrZvJwZ85rwuixhFr5pkTmaa")
    Wallet []ont.Wallet
    Acc  *ont.Account
    MulKey  = make([]keypair.PublicKey,0)
    acc  = make([]*ont.Account,0)
    MulSigAddr  common.Address
)

func main() {
    fmt.Println("===============================================================")
    fmt.Println("-----Ontology Transfer from MultiSign Address Tool started-----")
    fmt.Println("===============================================================")

    fmt.Println("Initial balance of address 1:")
    OntSdk = ont.NewOntologySdk()
    OntSdk.NewRpcClient().SetAddress("http://47.98.39.191:20336")
    balance1,_ := OntSdk.Native.Ont.BalanceOf(addr1)
    fmt.Println(balance1)

    fmt.Println("===============================================================")
    fmt.Println("Multi Keys:")
    for a := 1; a < 5; a++ {
        B := "wallets/path"
        A := "/wallet.dat"
        C := strconv.Itoa(a)
        Path := B + C + A
        wallet, err := ont.OpenWallet(Path)
        if err != nil {
            fmt.Println("account.Open error:", err)
            return
        }

        Acc, err = wallet.GetDefaultAccount(Passwd)
        if err != nil {
            fmt.Println("wallet.GetDefaultAccount error:", err)
            return
        }

        MulKey = append(MulKey, Acc.PublicKey)
        acc = append(acc, Acc)
    }
    fmt.Println(MulKey)

    fmt.Println("===============================================================")
    fmt.Println("Multi Sign Address:")
    MulSigAddr, _ = types.AddressFromMultiPubKeys(MulKey, 3)
    fmt.Println(MulSigAddr)

    fmt.Println("===============================================================")
    fmt.Println("Multi Sign Address balance:")
    balance,_ := OntSdk.Native.Ont.BalanceOf(MulSigAddr)
    fmt.Println(balance)

    States := []*ont2.State {{MulSigAddr, addr1, 1000000000}}
    tx, _ := OntSdk.Native.Ont.NewMultiTransferTransaction(GasPrice, GasLimit, States)

    for _, Acc = range acc {
        err := OntSdk.MultiSignToTransaction(tx, 3, MulKey, Acc)
        fmt.Println("MultiSignToTransaction error:", err)
    }

    fmt.Println("===============================================================")
    fmt.Println("Send Transaction result:")
    data,err := OntSdk.SendTransaction(tx)
    fmt.Println(data, err)

    fmt.Println("===============================================================")
    fmt.Println("End balance of address 1:")
    time.Sleep(time.Second * 10)
    balance1, _ = OntSdk.Native.Ont.BalanceOf(addr1)
    fmt.Println(balance1)

    fmt.Println("===============================================================")
    fmt.Println("------Ontology Transfer from MultiSign Address Tool Ended------")
    fmt.Println("===============================================================")
}