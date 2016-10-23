package context

import (
	"farm.e-pedion.com/repo/logger"
)

//Setup calls all provided setup functions and return all raised errors
func Setup(setupFuncs ...SetupFunc) []error {
	var errs []error
	for i, v := range setupFuncs {
		if err := v(); err != nil {
			logger.Warn("contex.Setup",
				logger.Int("index", i),
				logger.Struct("func", v),
				logger.Err(err),
			)
			errs = append(errs, err)
		}
	}
	return errs
}
