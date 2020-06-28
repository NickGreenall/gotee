package atomiser

// AtomData is used to encode parsed matches and subgroups.
type AtomData map[string]string

// Map is a helper function to perform a type cast.
func (data *AtomData) Map() map[string]string {
	return map[string]string(*data)
}
