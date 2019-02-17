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
	addFormulator(loader, ctd, common.MustParsePublicHash("2xASBuEWw6LcQGjYxeGZH9w1DUsEDt7fvUh8p3auxyN"), common.NewAddress(acg.Generate(), loader.ChainCoord(), 0))
	//addFormulator(loader, ctd, common.MustParsePublicHash("2VdGunZe8yZNm2mErqQqrFx2B7Mb4SBRPWviWnapahw"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
	/*
		addFormulator(loader, ctd, common.MustParsePublicHash("3eiovnNMgNCSkmxqwkjAabRTbNkkauMVk167Pgqon2Q"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
		addFormulator(loader, ctd, common.MustParsePublicHash("cNXbd7o43DkX48DaEy7hzuR6iy6DBxMAqNWmhxJLyA"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
		addFormulator(loader, ctd, common.MustParsePublicHash("3S7zbNCsAkHJns4Z3GP6RoQKcffHDxv8fPbk1tKD2Bb"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
		addFormulator(loader, ctd, common.MustParsePublicHash("39q6QQ9pfiP1yEAceCu11p5cmVhG8mHMiVayCD3UEa5"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
		addFormulator(loader, ctd, common.MustParsePublicHash("37pB69UiK7GX1sYcawoUq8c8yXS9WWbQnkmzoQjUmZB"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
		addFormulator(loader, ctd, common.MustParsePublicHash("2r9mQmdfvK62ELWezK8tUvDztettkUkEGrvMWUXL7D"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
		addFormulator(loader, ctd, common.MustParsePublicHash("2CQBhmtferf2qWDjqSnEE3f1ECimj4Lck2CxndgqEVq"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
		addFormulator(loader, ctd, common.MustParsePublicHash("4D5m6ssnsf3NxJmqKg7PpwoyG2PdMNPAuQjpB8ZKjDo"), common.NewAddress(acg.Generate(), st.ChainCoord(), 0))
	*/
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
