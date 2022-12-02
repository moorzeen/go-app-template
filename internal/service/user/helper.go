package user

func passComplexity(pass string) error {
	if len([]rune(pass)) < 8 {
		return ErrWeakPassword
	}
	return nil
}
