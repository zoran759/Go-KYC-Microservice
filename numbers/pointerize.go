package numbers

func PointerizeInt(i int) *int {
	if i == 0 {
		return nil
	}

	return &i
}
