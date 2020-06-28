package atomiser

func (atom AtomData) MockMarshal() (interface{}, error) {
	return atom, nil
}

func (atom AtomData) MockUnmarshal(v interface{}) error {
	target, ok := v.(*AtomData)
	if !ok {
		return nil
	}
	*target = atom
	return nil
}
