package vos

type (
	TransactionID string
	AccountID     string
)

// String returns transaction id as string
func (t TransactionID) String() string {
	return string(t)
}

// String returns account id as string
func (a AccountID) String() string {
	return string(a)
}
