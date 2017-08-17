package nodes

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gebv/as_gifts/errors"
)

func NewInMemoryStoreNodes() *InMemoryStoreNodes {
	return &InMemoryStoreNodes{
		List: make(map[string]Node),
	}
}

type InMemoryStoreNodes struct {
	sync.Mutex

	List map[string]Node
}

func (s *InMemoryStoreNodes) Find(chatID string) (*Node, error) {
	s.Lock()
	defer s.Unlock()

	obj, exists := s.List[chatID]
	if !exists {
		return nil, errors.ErrNotFound
	}
	return &obj, nil
}

func (s *InMemoryStoreNodes) LoadFromToml(dat []byte) error {
	obj := struct {
		Nodes []Node `toml:"nodes"`
	}{}
	if err := toml.Unmarshal(dat, &obj); err != nil {
		return err
	}
	for _, node := range obj.Nodes {
		s.List[node.ID] = node
	}
	return nil
}
