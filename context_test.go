package context

import (
	"errors"
	"farm.e-pedion.com/repo/logger"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	os.Args = append(os.Args, "-ecf", "./test/etc/context/context.yaml")
	logger.Info("context_test.init")
}

func TestSetupAll(t *testing.T) {
	errs := SetupAll(
		func() error { return nil },
		func() error { return nil },
		func() error { return nil },
		func() error { return nil },
	)

	assert.Nil(t, errs)
	assert.Empty(t, errs)
}

func TestSetupAllErr(t *testing.T) {
	var errs []error
	assert.NotPanics(t, func() {
		errs = SetupAll(
			func() error { return errors.New("context_test.TestSetupErrMock1") },
			func() error { return errors.New("context_test.TestSetupErrMock2") },
			func() error { return errors.New("context_test.TestSetupErrMock3") },
			func() error { return errors.New("context_test.TestSetupErrMock4") },
		)
	})

	assert.NotNil(t, errs)
	assert.NotEmpty(t, errs)
	assert.Equal(t, 4, len(errs))
}

func TestSetup(t *testing.T) {
	errs := Setup(
		func() error { return nil },
		func() error { return nil },
		func() error { return nil },
		func() error { return nil },
	)

	assert.Nil(t, errs)
	assert.Empty(t, errs)
}

func TestSetupErr(t *testing.T) {
	firstErrMsg := "context_test.TestSetupErrMock1"
	var err error
	assert.NotPanics(t, func() {
		err = Setup(
			func() error { return errors.New(firstErrMsg) },
			func() error { return errors.New("context_test.TestSetupErrMock2") },
			func() error { return errors.New("context_test.TestSetupErrMock3") },
			func() error { return errors.New("context_test.TestSetupErrMock4") },
		)
	})

	assert.NotNil(t, err)
	assert.Equal(t, firstErrMsg, err.Error())
}
