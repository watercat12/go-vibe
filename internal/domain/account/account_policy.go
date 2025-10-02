package account

type AccountPolicy struct {}

func NewAccountPolicy() *AccountPolicy {
	return &AccountPolicy{}
}

func (p *AccountPolicy) LimitSavingAccount(count int) bool {
	return count >= 5
}