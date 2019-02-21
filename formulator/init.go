package main

import (
	"log"

	"git.fleta.io/fleta/common"
	"git.fleta.io/fleta/core/account"
	"git.fleta.io/fleta/core/amount"
	"git.fleta.io/fleta/core/consensus"
	"git.fleta.io/fleta/core/data"
	"git.fleta.io/fleta/core/transaction"
	"git.fleta.io/fleta/extension/account_def"

	_ "git.fleta.io/fleta/extension/account_tx"
	_ "git.fleta.io/fleta/extension/utxo_tx"
	_ "git.fleta.io/fleta/solidity"
)

// consts
const (
	BlockchainVersion = 1
)

func initChainComponent(act *data.Accounter, tran *data.Transactor) error {
	// transaction_type transaction types
	const (
		// FLETA Transactions
		TransferTransctionType              = transaction.Type(10)
		WithdrawTransctionType              = transaction.Type(18)
		BurnTransctionType                  = transaction.Type(19)
		CreateAccountTransctionType         = transaction.Type(20)
		CreateMultiSigAccountTransctionType = transaction.Type(21)
		// UTXO Transactions
		AssignTransctionType      = transaction.Type(30)
		DepositTransctionType     = transaction.Type(38)
		OpenAccountTransctionType = transaction.Type(41)
		// Formulation Transactions
		CreateFormulationTransctionType = transaction.Type(60)
		RevokeFormulationTransctionType = transaction.Type(61)
		// Solidity Transactions
		SolidityCreateContractType = transaction.Type(70)
		SolidityCallContractType   = transaction.Type(71)
	)

	// account_type account types
	const (
		// FLTEA Accounts
		SingleAccountType   = account.Type(10)
		MultiSigAccountType = account.Type(11)
		LockedAccountType   = account.Type(19)
		// Formulation Accounts
		FormulationAccountType = account.Type(60)
		// Solidity Accounts
		SolidityAccount = account.Type(70)
	)

	type txFee struct {
		Type transaction.Type
		Fee  *amount.Amount
	}

	TxFeeTable := map[string]*txFee{
		"fleta.CreateAccount":         &txFee{CreateAccountTransctionType, amount.COIN.MulC(10)},
		"fleta.CreateMultiSigAccount": &txFee{CreateMultiSigAccountTransctionType, amount.COIN.MulC(10)},
		"fleta.Transfer":              &txFee{TransferTransctionType, amount.COIN.DivC(10)},
		"fleta.Withdraw":              &txFee{WithdrawTransctionType, amount.COIN.DivC(10)},
		"fleta.Burn":                  &txFee{BurnTransctionType, amount.COIN.DivC(10)},
		"fleta.Assign":                &txFee{AssignTransctionType, amount.COIN.DivC(2)},
		"fleta.Deposit":               &txFee{DepositTransctionType, amount.COIN.DivC(2)},
		"fleta.OpenAccount":           &txFee{OpenAccountTransctionType, amount.COIN.MulC(10)},
		"consensus.CreateFormulation": &txFee{CreateFormulationTransctionType, amount.COIN.MulC(50000)},
		"consensus.RevokeFormulation": &txFee{RevokeFormulationTransctionType, amount.COIN.DivC(10)},
		"solidity.CreateContract":     &txFee{SolidityCreateContractType, amount.COIN.MulC(10)},
		"solidity.CallContract":       &txFee{SolidityCallContractType, amount.COIN.DivC(10)},
	}
	for name, item := range TxFeeTable {
		if err := tran.RegisterType(name, item.Type, item.Fee); err != nil {
			log.Println(name, item, err)
			return err
		}
	}

	AccTable := map[string]account.Type{
		"fleta.SingleAccount":          SingleAccountType,
		"fleta.MultiSigAccount":        MultiSigAccountType,
		"fleta.LockedAccount":          LockedAccountType,
		"consensus.FormulationAccount": FormulationAccountType,
		"solidity.ContractAccount":     SolidityAccount,
	}
	for name, t := range AccTable {
		if err := act.RegisterType(name, t); err != nil {
			log.Println(name, t, err)
			return err
		}
	}
	return nil
}

func initGenesisContextData(act *data.Accounter, tran *data.Transactor) (*data.ContextData, error) {
	loader := data.NewEmptyLoader(act.ChainCoord(), act, tran)
	ctd := data.NewContextData(loader, nil)

	acg := &accCoordGenerator{}
	addSingleAccount(loader, ctd, common.MustParsePublicHash("3Zmc4bGPP7TuMYxZZdUhA9kVjukdsE2S8Xpbj4Laovv"), common.NewAddress(acg.Generate(), loader.ChainCoord(), 0))
	addFormulator(loader, ctd, common.MustParsePublicHash("gDGAcf9V9i8oWLTeayoKC8bdAooNVaFnAeQKy4CsUB"), common.MustParseAddress("3CUsUpvEK"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4m6XsJbq6EFb5bqhZuKFc99SmF86ymcLcRPwrWyToHQ"), common.MustParseAddress("5PxjxeqTd"))
	addFormulator(loader, ctd, common.MustParsePublicHash("o1rVoXHFuz5EtwLwCLcrmHpqPdugAnWHEVVMtnCb32"), common.MustParseAddress("7bScSUkgw"))
	addFormulator(loader, ctd, common.MustParsePublicHash("47NZ8oadY4dCAM3ZrGFrENPn99L1SLSqzpR4DFPUpk5"), common.MustParseAddress("9nvUvJfvF"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4TaHVFSzcrNPktRiNdpPitoUgLXtZzrVmkxE3GmcYjG"), common.MustParseAddress("BzQMQ8b9Z"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2wqsb4J47T4JkNUp1Bma1HkjpCyei7sZinLmNprpdtY"), common.MustParseAddress("EBtDsxWNs"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2a1CirwCHSYYpLqpbi1b7Rpr4BAJZvydbDA1bGjJ7FG"), common.MustParseAddress("GPN6MnRcB"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2KnMHH973ZLicENxcsJbARdeTUiYZmN3WnBzbZqvvEx"), common.MustParseAddress("JaqxqcLqV"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4fyTmraz8x3NKWnj4nWgPWKy8qCBF1hyqVJQeyupHAe"), common.MustParseAddress("LnKqKSG4o"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2V1zboMnJbJdeLvRBRFVPvVqs8CCmjxToBpGJSNScu2"), common.MustParseAddress("NyohoGBJ7"))
	addFormulator(loader, ctd, common.MustParsePublicHash("3pEYkEgXoPUm4vdcGBXP46q1BpMj215uVQdAg6P4g74"), common.MustParseAddress("RBHaH66XR"))
	addFormulator(loader, ctd, common.MustParsePublicHash("rsUoPRfVgXJFuV6wYcy4M4kntvr3tooeXzcRhrjBq6"), common.MustParseAddress("TNmSkv1kj"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4UMYzaBeXEKcm6hnDDEMqYRR5NLwGndCLksryVj98Fw"), common.MustParseAddress("VaFKEjvz3"))
	addFormulator(loader, ctd, common.MustParsePublicHash("3h2Lt2uYFMqVQKFgKszLJzwaLhQ5kt1nMcg8M758aLh"), common.MustParseAddress("XmjBiZrDM"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4NkvvfPdHHvpo9YTkAQBrGxpnnML2pVRXHdLgzB2EYe"), common.MustParseAddress("ZyD4CPmSf"))
	addFormulator(loader, ctd, common.MustParsePublicHash("3ae9sCuM75vAheVLNp3DjQqDiD3TaxY5HYduHvsgzYZ"), common.MustParseAddress("cAgvgDgfy"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2bR5L2ZSqKLUFQzdhzWV6e4BUupHPGDFtnZUNrZBZbZ"), common.MustParseAddress("eNAoA3buH"))
	addFormulator(loader, ctd, common.MustParsePublicHash("BPqzvcrYi364mm6GyraHHqJHrvEfqjwo1jEC8crTxZ"), common.MustParseAddress("gZefdsX8b"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2vtYXNUAtBtt4fF6DEbVKNc7bGhA7yBbatTA6Ye9kMT"), common.MustParseAddress("im8Y7hSMu"))
	addFormulator(loader, ctd, common.MustParsePublicHash("42TUBLNb1natk7s7qsHNqxHwn7Pb3pNmTfTnd1sDQnb"), common.MustParseAddress("kxcQbXMbD"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2yng1DwwBqMixjCnjx6Pdf9o5AkgEzkumxJySr8Qe6C"), common.MustParseAddress("oA6H5MGpX"))
	addFormulator(loader, ctd, common.MustParsePublicHash("3PNrAwb7FrvKeB1hCxYADwNxqWuYmaqoc8E8VjdBC"), common.MustParseAddress("qMa9ZBC3q"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2eZAofvjk5AHUpaUyC7EDx3K8KAHUQNXMynHG7ZYFfn"), common.MustParseAddress("sZ42317H9"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4QT4FGpoaFkPiRaZQCKDfrANWJ6EAqavqkQfGr6g4oG"), common.MustParseAddress("ukXtWq2WT"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2nPZHDpFavW2VjnZGs7ZeQyFM19y517ZTQaTgqe3G69"), common.MustParseAddress("wx1kzewjm"))
	addFormulator(loader, ctd, common.MustParsePublicHash("bB88uMhpM4vjUHpV5WZqfQBh4kyi6wnnKCtVF4AE2D"), common.MustParseAddress("z9VdUUry5"))
	addFormulator(loader, ctd, common.MustParsePublicHash("2ZLEXwQ9pqvaATFttkkNWY2CGDHdJFa5V3GNapKeqtx"), common.MustParseAddress("22LyVxJnCP"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4M2KFgmWSKu8JyjhkmVJ8U4hjtn9MX4rsch4ZoE1i32"), common.MustParseAddress("24YTNS8hRh"))
	addFormulator(loader, ctd, common.MustParsePublicHash("XG9nFJsdMpo6D6wYxYSyH5zAtnvsMjySFHp1XjCouY"), common.MustParseAddress("26jwEuxcf1"))
	addFormulator(loader, ctd, common.MustParsePublicHash("3uW4bb1kAx35ndj4ZVLMF8xWYercS2RfP7moxZvUm8Y"), common.MustParseAddress("28wR7PnXtK"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4mY5G1BZuZaeHR5cH1K4sUNmccPa11JkHtjv5ctde3K"), common.MustParseAddress("2B8tyscT7d"))
	addFormulator(loader, ctd, common.MustParsePublicHash("3oocpeXtqUZeaut1A71fbCMBQefMFMCBt2BpamNZfA9"), common.MustParseAddress("2DLNrMSNLw"))
	addFormulator(loader, ctd, common.MustParsePublicHash("4wknRQ86rTcN1cQbXZfbCMkqXcS1FsYG8ihAYFhmxF"), common.MustParseAddress("2FXriqGHaF"))
	addFormulator(loader, ctd, common.MustParsePublicHash("3mT9SNvGscpwmDjHnojnVysd9pXUvg1fenVyiBFYTDs"), common.MustParseAddress("2HjLbK6CoZ"))
	addFormulator(loader, ctd, common.MustParsePublicHash("24zn1BgQBmMD8dWap9XbBHdZAivDppVhnYxzZ4ftZw4"), common.MustParseAddress("2KvpTnv82s"))
	return ctd, nil
}

func addSingleAccount(loader data.Loader, ctd *data.ContextData, KeyHash common.PublicHash, addr common.Address) {
	a, err := loader.Accounter().NewByTypeName("fleta.SingleAccount")
	if err != nil {
		panic(err)
	}
	acc := a.(*account_def.SingleAccount)
	acc.Address_ = addr
	acc.KeyHash = KeyHash
	ctd.CreatedAccountMap[acc.Address_] = acc
	balance := account.NewBalance()
	balance.AddBalance(loader.ChainCoord(), amount.NewCoinAmount(10000000000, 0))
	ctd.AccountBalanceMap[acc.Address_] = balance
}

func addFormulator(loader data.Loader, ctd *data.ContextData, KeyHash common.PublicHash, addr common.Address) {
	a, err := loader.Accounter().NewByTypeName("consensus.FormulationAccount")
	if err != nil {
		panic(err)
	}
	acc := a.(*consensus.FormulationAccount)
	acc.Address_ = addr
	acc.KeyHash = KeyHash
	ctd.CreatedAccountMap[acc.Address_] = acc
}

type accCoordGenerator struct {
	idx uint16
}

func (acg *accCoordGenerator) Generate() *common.Coordinate {
	coord := common.NewCoordinate(0, acg.idx)
	acg.idx++
	return coord
}
