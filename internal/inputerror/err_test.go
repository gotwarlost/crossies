package inputerror_test

import (
	"testing"

	"github.com/gotwarlost/crossies/internal/inputerror"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	e := inputerror.New("foobar")
	assert.Equal(t, "foobar", e.Error())
	assert.True(t, inputerror.IsInputError(e))
	e2 := errors.Wrap(e, "barbaz")
	assert.Equal(t, "barbaz: foobar", e2.Error())
	assert.True(t, inputerror.IsInputError(e2))
}
