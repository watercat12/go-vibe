package bank_link

type BankLinkError struct {
	Message string
}

func (e *BankLinkError) Error() string {
	return e.Message
}