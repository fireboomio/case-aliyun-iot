package dict

import (
	"custom-go/generated"
	"custom-go/operation/dict/item/upsertMany"
	"custom-go/pkg/types"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"sync"
)

type SortedMap[K comparable, V any] struct {
	dictCode   string
	kvFetcher  func(k, v string) (K, V)
	keySorter  func(K, K) bool
	keyMatcher func(K, K) bool
	defaultVal V
	sync.Mutex

	keys []K
	data map[K]V
}

func NewSortedMap[K comparable, V any](dictCode string,
	kvFetcher func(k, v string) (K, V), keySorter func(K, K) bool,
	keyMatcher func(K, K) bool, defaultVal V) *SortedMap[K, V] {
	_sortedMap := &SortedMap[K, V]{
		dictCode:   dictCode,
		kvFetcher:  kvFetcher,
		keySorter:  keySorter,
		keyMatcher: keyMatcher,
		defaultVal: defaultVal,
	}
	_sortedMap.refresh(types.NewEmptyInternalClient())
	upsertMany.Subscribe(dictCode, _sortedMap.refresh)
	return _sortedMap
}

func (m *SortedMap[K, V]) refresh(client *types.InternalClient) {
	m.Lock()
	defer m.Unlock()
	dictFindInput := generated.Dict__item__findManyInternalInput{DictCode: m.dictCode, Enabled: true}
	dictFindResp, _ := generated.Dict__item__findMany.Execute(dictFindInput, client)
	_dictData := make(map[K]V, len(dictFindResp.Data))
	for _, item := range dictFindResp.Data {
		k, v := m.kvFetcher(item.Key, item.Value)
		_dictData[k] = v
	}
	_dictKeys := maps.Keys(_dictData)
	if m.keySorter != nil {
		slices.SortFunc(_dictKeys, m.keySorter)
	}
	m.keys, m.data = _dictKeys, _dictData
}

func (m *SortedMap[K, V]) Search(key K) (v V, k K) {
	m.Lock()
	defer m.Unlock()
	for _, _k := range m.keys {
		if m.keyMatcher(key, _k) {
			return m.data[_k], _k
		}
	}
	v = m.defaultVal
	return
}
