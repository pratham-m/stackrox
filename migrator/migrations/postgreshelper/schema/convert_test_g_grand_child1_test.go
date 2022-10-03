// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTestGGrandChild1Serialization(t *testing.T) {
	obj := &storage.TestGGrandChild1{}
	assert.NoError(t, testutils.FullInit(obj, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	m, err := ConvertTestGGrandChild1FromProto(obj)
	assert.NoError(t, err)
	conv, err := ConvertTestGGrandChild1ToProto(m)
	assert.NoError(t, err)
	assert.Equal(t, obj, conv)
}
