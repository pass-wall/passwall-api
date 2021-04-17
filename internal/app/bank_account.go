package app

import (
	"github.com/passwall/passwall-server/internal/storage"
	"github.com/passwall/passwall-server/model"
)

// CreateBankAccount creates a new bank account and saves it to the store
func CreateBankAccount(s storage.Store, dto *model.BankAccountDTO, schema string) (*model.BankAccount, error) {
	createdBankAccount, err := s.BankAccounts().Save(EncryptModel(model.ToBankAccount(dto)).(*model.BankAccount), schema)
	if err != nil {
		return nil, err
	}

	return createdBankAccount, nil
}

// UpdateBankAccount updates the account with the dto and applies the changes in the store
func UpdateBankAccount(s storage.Store, bankAccount *model.BankAccount, dto *model.BankAccountDTO, schema string) (*model.BankAccount, error) {
	encModel := EncryptModel(model.ToBankAccount(dto)).(*model.BankAccount)

	bankAccount.BankName = encModel.BankName
	bankAccount.BankCode = encModel.BankCode
	bankAccount.AccountName = encModel.AccountName
	bankAccount.AccountNumber = encModel.AccountNumber
	bankAccount.IBAN = encModel.IBAN
	bankAccount.Currency = encModel.Currency
	bankAccount.Password = encModel.Password

	updatedBankAccount, err := s.BankAccounts().Save(bankAccount, schema)
	if err != nil {
		return nil, err
	}

	return updatedBankAccount, nil
}
