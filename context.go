package context

//Setup calls all provided setup functions and return all raised errors
func Setup(setupFuncs ...SetupFunc) []error {
	var errs []error
	for _, v := range setupFuncs {
		if err := v(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
