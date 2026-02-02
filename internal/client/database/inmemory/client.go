package inmemory


type InMemoryClient struct {
	db *InMemoryDB
}

func NewInMemoryClient() *InMemoryClient {
	return &InMemoryClient{
		db: NewInMemoryDB(),
	}
}

func (c *InMemoryClient) CreateTable(name string) error {
	newTable := NewInMemoryDBTable(name)

	return c.db.AddTable(newTable)
}

func (c *InMemoryClient) GetTable(name string) (IInMemoryDBTable, error) {
	return c.db.GetTable(name)
}	
