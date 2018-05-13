package config

type Transaction struct {
	*Query
	Resp chan int
}

func NewTransaction(query *Query) *Transaction {
	return &Transaction{
		Query: query,
		Resp:  make(chan int, 1),
	}
}
