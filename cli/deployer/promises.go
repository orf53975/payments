package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mysteriumnetwork/payments/mysttoken"
	"github.com/mysteriumnetwork/payments/promises"
)

func DeployPromises(transactor *bind.TransactOpts, client bind.ContractBackend, erc20tokenAddress common.Address, registrationFee int64) error {
	erc20, err := mysttoken.NewERC20(erc20tokenAddress, client)
	if err != nil {
		return err
	}

	deployerBalance, err := erc20.BalanceOf(&bind.CallOpts{}, transactor.From)
	if err != nil {
		return err
	}
	fmt.Println("Deployer balance of erc20 is: " + deployerBalance.String())

	pc, err := promises.DeployPromiseClearer(transactor, erc20tokenAddress, registrationFee, client)
	if err != nil {
		return err
	}
	fmt.Println("Deployed contract address: " + pc.Address.String())
	return nil
}
