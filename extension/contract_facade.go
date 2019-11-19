package extension

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/extension/extensionContracts"
	"math/big"
)

type ManagementContractFacade interface {
	Transactor(managementAddress common.Address) (*extensionContracts.ContractExtenderTransactor, error)
	Caller(managementAddress common.Address) (*extensionContracts.ContractExtenderCaller, error)
	Deploy(args *bind.TransactOpts, toExtend common.Address, voters []common.Address, recipientHash string, uuid string) (*types.Transaction, error)

	GetAllVoters(addressToVoteOn common.Address) ([]common.Address, error)
}

type EthclientManagementContractFacade struct {
	client *ethclient.Client
}

func NewManagementContractFacade(client *ethclient.Client) ManagementContractFacade {
	return EthclientManagementContractFacade{client: client}
}

func (facade EthclientManagementContractFacade) Transactor(managementAddress common.Address) (*extensionContracts.ContractExtenderTransactor, error) {
	return extensionContracts.NewContractExtenderTransactor(managementAddress, facade.client)
}

func (facade EthclientManagementContractFacade) Caller(managementAddress common.Address) (*extensionContracts.ContractExtenderCaller, error) {
	return extensionContracts.NewContractExtenderCaller(managementAddress, facade.client)
}

func (facade EthclientManagementContractFacade) Deploy(args *bind.TransactOpts, toExtend common.Address, voters []common.Address, recipientHash string, uuid string) (*types.Transaction, error) {
	_, tx, _, err := extensionContracts.DeployContractExtender(args, facade.client, toExtend, voters, recipientHash, uuid)
	return tx, err
}

/////////

func (facade EthclientManagementContractFacade) GetAllVoters(addressToVoteOn common.Address) ([]common.Address, error) {
	caller, err := extensionContracts.NewContractExtenderCaller(addressToVoteOn, facade.client)
	if err != nil {
		return nil, err
	}
	numberOfVoters, err := caller.TotalNumberOfVoters(nil)
	if err != nil {
		return nil, err
	}
	var i int64
	var voters []common.Address
	for i = 0; i < numberOfVoters.Int64(); i++ {
		voter, err := caller.WalletAddressesToVote(nil, big.NewInt(i))
		if err != nil {
			return nil, err
		}
		voters = append(voters, voter)
	}
	return voters, nil
}
