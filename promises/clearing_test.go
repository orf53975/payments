package promises

import (
	"math/big"
	"testing"
	"time"

	"github.com/mysteriumnetwork/payments/contracts/abigen"
	"github.com/mysteriumnetwork/payments/mysttoken"
	"github.com/mysteriumnetwork/payments/test_utils"
	"github.com/stretchr/testify/assert"
)

var abiList, _ = test_utils.ParseAbis(test_utils.AbiMap{
	"MystToken": {
		mysttoken.MystTokenABI,
		mysttoken.MystTokenBin,
	},
	"IdentityPromises": {
		abigen.IdentityPromisesABI,
		abigen.IdentityPromisesBin,
	},
})

func TestPromiseClearingEmitsClearedEvent(t *testing.T) {
	backend := test_utils.LoggingBackend(test_utils.NewSimulatedBackend(test_utils.Deployer.Address, 10000000000), abiList)

	mystErc20, err := mysttoken.DeployMystERC20(test_utils.Deployer.Transactor, 1000000, backend)
	assert.NoError(t, err)

	clearing, err := DeployPromiseClearer(test_utils.Deployer.Transactor, mystErc20.Address, 1000, backend)
	assert.NoError(t, err)
	backend.Commit()

	_, err = mystErc20.Approve(clearing.Address, big.NewInt(3000))
	assert.NoError(t, err)
	backend.Commit()

	events := make(chan *abigen.IdentityPromisesPromiseCleared, 1)
	sub, err := clearing.BindForEvents(events)
	assert.NoError(t, err)

	payer, err := test_utils.NewMystIdentity()
	assert.NoError(t, err)

	receiver, err := test_utils.NewMystIdentity()
	assert.NoError(t, err)

	err = clearing.RegisterIdentities(*payer, *receiver)
	assert.NoError(t, err)
	backend.Commit()

	_, err = clearing.TopUp(payer.Address, big.NewInt(1000))
	assert.NoError(t, err)
	backend.Commit()

	promise := Promise{
		Extra:    EmptyExtra{},
		Receiver: receiver.Address,
		SeqNo:    1,
		Amount:   100,
	}

	issuedPromise, err := SignByPayer(&promise, payer)
	assert.NoError(t, err)
	receivedPromise, err := SignByReceiver(issuedPromise, receiver)
	assert.NoError(t, err)

	err = clearing.ClearReceivedPromise(receivedPromise)
	assert.NoError(t, err)
	backend.Commit()

	lastSeqNo, err := clearing.LastClearedPromise(payer.Address, receiver.Address)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), lastSeqNo)

	balanceOfSender, err := clearing.Balances(payer.Address)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(900), balanceOfSender)

	balanceOfReceiver, err := clearing.Balances(receiver.Address)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(100), balanceOfReceiver)

	select {
	case event := <-events:
		assert.Equal(t, big.NewInt(1), event.SeqNo)
		assert.Equal(t, big.NewInt(100), event.Amount)
		assert.Equal(t, payer.Address, event.From)
		assert.Equal(t, receiver.Address, event.To)
	case err := <-sub.Err():
		assert.NoError(t, err)
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Event from contract expected")
	}

	sub.Unsubscribe()
}
