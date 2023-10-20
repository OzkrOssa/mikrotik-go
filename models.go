package mikrotik

type Secret struct {
	Name     string
	CallerID string
	Profile  string
	Comment  string
}

type AddressList struct {
	ID           string
	Address      string
	Comment      string
	CreationTime string
	List         string
	Status       string
}
