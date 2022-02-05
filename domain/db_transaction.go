package domain

type DBTransactionRepository interface {
	Transaction(func(interface{}) error) error
}
