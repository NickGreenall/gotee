package atomiser

type AtomData map[string][]byte

func (data *AtomData) Map() map[string][]byte {
	return map[string][]byte(*data)
}
