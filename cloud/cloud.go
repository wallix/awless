package cloud

type Cloud interface {
	Identifier
}

type Identifier interface {
	GetAccountId() (string, error)
	GetUserId() (string, error)
}

var (
	Current Cloud
)
