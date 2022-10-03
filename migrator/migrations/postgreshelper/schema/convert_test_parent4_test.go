// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTestParent4Serialization(t *testing.T) {
	obj := &storage.TestParent4{}
	assert.NoError(t, testutils.FullInit(obj, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	m, err := ConvertTestParent4FromProto(obj)
	assert.NoError(t, err)
	conv, err := ConvertTestParent4ToProto(m)
	assert.NoError(t, err)
	assert.Equal(t, obj, conv)
}
