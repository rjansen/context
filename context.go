package context

import (
	"farm.e-pedion.com/repo/logger"
)

//SetupSilent calls all provided setup functions and return all raised errors
func SetupSilent(setupFuncs ...SetupFunc) []error {
	var errs []error
	for i, v := range setupFuncs {
		if err := v(); err != nil {
			logger.Warn("contex.SetupSilent",
				logger.Int("index", i),
				logger.Struct("func", v),
				logger.Err(err),
			)
			errs = append(errs, err)
		}
	}
	return errs
}

//Setup calls the provided setup functions and return at the first raised error
func Setup(setupFuncs ...SetupFunc) error {
	for i, v := range setupFuncs {
		if err := v(); err != nil {
			logger.Error("contex.Setup",
				logger.Int("index", i),
				logger.Struct("func", v),
				logger.Err(err),
			)
			return err
		}
	}
	return nil
}
