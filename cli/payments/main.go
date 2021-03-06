package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mysteriumnetwork/payments/cli/helpers"
	"github.com/mysteriumnetwork/payments/contracts/abigen"
)

var paymentsContract = flag.String("payments.contract", "", "Address of payments contract")
var identity = flag.String("payments.identity", "", "Identity for balance checking")
var gethUrl = flag.String("geth.url", "", "URL value of started geth to connect")

func main() {
	flag.Parse()

	client, syncCompleted, err := helpers.LookupBackend(*gethUrl)
	checkError(err)
	<-syncCompleted

	contractCaller, err := abigen.NewIdentityPromisesCaller(common.HexToAddress(*paymentsContract), client)
	checkError(err)

	paymentsContract := abigen.IdentityPromisesCallerSession{
		Contract: contractCaller,
		CallOpts: bind.CallOpts{},
	}

	registered, err := paymentsContract.IsRegistered(common.HexToAddress(*identity))
	checkError(err)
	fmt.Println("Identity: ", *identity, "registration status: ", registered)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}
