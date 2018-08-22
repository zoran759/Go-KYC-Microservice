package strings

func Pointerize(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
