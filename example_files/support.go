package src

func returns2Values() (int, error) {
	return 0, nil
}

func returns3Values() (int, string, error) {
	return 0, "", nil
}

func returnsBool() (int, bool) {
	return 0, true
}

func SingleErrReturn() error {
	return nil
}

type ErrStruct struct {
	field1 int
	field2 error
}
