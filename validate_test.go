package ambient_test

import (
	"context"
	"testing"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	// Allow lowercase name and proper version.
	p := mock.NewPlugin("test", "1.0.0")
	err := ambient.Validate(context.Background(), p)
	assert.NoError(t, err)

	// Allow single letter name.
	p = mock.NewPlugin("t", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.NoError(t, err)

	// Don't allow starting with a number name.
	p = mock.NewPlugin("7test", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow underscore in name.
	p = mock.NewPlugin("te_st", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow space in name.
	p = mock.NewPlugin("te st", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow blank name.
	p = mock.NewPlugin("", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow name: plugin.
	p = mock.NewPlugin("plugin", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow start with capital letter.
	p = mock.NewPlugin("Test", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow any capital letter.
	p = mock.NewPlugin("tesT", "1.0.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow bad version.
	p = mock.NewPlugin("test", "1.0.")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Allow version with 2 numbers.
	p = mock.NewPlugin("test", "1.0")
	err = ambient.Validate(context.Background(), p)
	assert.NoError(t, err)

	// Don't allow bad version.
	p = mock.NewPlugin("test", "1.")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Allow version with 1 number.
	p = mock.NewPlugin("test", "1")
	err = ambient.Validate(context.Background(), p)
	assert.NoError(t, err)

	// Don't allow bad version.
	p = mock.NewPlugin("test", "-1")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow bad version.
	p = mock.NewPlugin("test", "1.0-a")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Don't allow bad version with v passed in.
	p = mock.NewPlugin("test", "v1.0")
	err = ambient.Validate(context.Background(), p)
	assert.Error(t, err)

	// Allow version with pre-release.
	p = mock.NewPlugin("test", "1.0.0-alpha")
	err = ambient.Validate(context.Background(), p)
	assert.NoError(t, err)

	// Allow version with pre-release and build
	p = mock.NewPlugin("test", "1.0.0-alpha+001")
	err = ambient.Validate(context.Background(), p)
	assert.NoError(t, err)

	// Allow version with build
	p = mock.NewPlugin("test", "1.0.0+001")
	err = ambient.Validate(context.Background(), p)
	assert.NoError(t, err)

}
