package inmemory

import (
	"fmt"
	"github.com/Peqchji/go-inbound-adapter-benchmark/pkg"
	"sync"
)

type InMemoryDB struct {
	name       string
	mu         sync.RWMutex
	tableStore map[string]IInMemoryDBTable
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		name:       "InMemoryDB",
		tableStore: make(map[string]IInMemoryDBTable),
	}
}

func (r *InMemoryDB) AddTable(table IInMemoryDBTable) error {
	if table == nil {
		return fmt.Errorf("%s: %w", r.name, ErrCreateTableFailed)
	}

	name := table.Name()
	if _, ok := r.tableStore[name]; ok {
		return fmt.Errorf("%s: %w", r.name, ErrNameTableCollision)
	}

	r.tableStore[name] = table

	return nil
}

func (r *InMemoryDB) GetTable(name string) (IInMemoryDBTable, error) {
	if table, ok := r.tableStore[name]; ok {
		return table, nil
	}

	return nil, fmt.Errorf("%s: %w <%s>", r.name, ErrNotFoundTable, name)
}

//---------------------------------------------------------------------------//

type IInMemoryDBTable interface {
	Name() string
	GetById(id string) pkg.Result[pkg.Record]
	Save(data pkg.Record) error
	GetAll() pkg.Result[[]pkg.Record]
}

var _ IInMemoryDBTable = &InMemoryDBTable{}

type InMemoryDBTable struct {
	name string
	mu   sync.RWMutex
	data map[string]pkg.Record
}

func NewInMemoryDBTable(name string) *InMemoryDBTable {
	return &InMemoryDBTable{
		name: fmt.Sprintf("InMemoryDBTable_%s", name),
		data: make(map[string]pkg.Record),
	}
}

func (t *InMemoryDBTable) Name() string {
	return t.name
}

func (t *InMemoryDBTable) GetById(id string) pkg.Result[pkg.Record] {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if record, ok := t.data[id]; ok {
		return pkg.Result[pkg.Record]{
			Res: record,
			Err: nil,
		}
	}

	return pkg.Result[pkg.Record]{
		Res: nil,
		Err: ErrNotFoundRecord,
	}
}

func (t *InMemoryDBTable) Save(data pkg.Record) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrSaveError
		}
	}()

	t.mu.Lock()
	defer t.mu.Unlock()

	dataId := data.GetID()
	t.data[dataId] = data

	return err
}

func (t *InMemoryDBTable) GetAll() pkg.Result[[]pkg.Record] {
	t.mu.RLock()
	defer t.mu.RUnlock()

	amount := len(t.data)
	if amount == 0 {
		return pkg.Result[[]pkg.Record]{
			Res: nil,
			Err: ErrNotFoundRecord,
		}
	}

	i := 0
	records := make([]pkg.Record, amount)
	for _, rec := range t.data {
		records[i] = rec
		i += 1
	}

	return pkg.Result[[]pkg.Record]{
		Res: records,
		Err: nil,
	}
}
