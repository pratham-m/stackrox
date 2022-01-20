// Code generated by MockGen. DO NOT EDIT.
// Source: indexer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

// MockIndexer is a mock of Indexer interface.
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer.
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance.
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// AddComponentCVEEdge mocks base method.
func (m *MockIndexer) AddComponentCVEEdge(componentcveedge *storage.ComponentCVEEdge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComponentCVEEdge", componentcveedge)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddComponentCVEEdge indicates an expected call of AddComponentCVEEdge.
func (mr *MockIndexerMockRecorder) AddComponentCVEEdge(componentcveedge interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComponentCVEEdge", reflect.TypeOf((*MockIndexer)(nil).AddComponentCVEEdge), componentcveedge)
}

// AddComponentCVEEdges mocks base method.
func (m *MockIndexer) AddComponentCVEEdges(componentcveedges []*storage.ComponentCVEEdge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComponentCVEEdges", componentcveedges)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddComponentCVEEdges indicates an expected call of AddComponentCVEEdges.
func (mr *MockIndexerMockRecorder) AddComponentCVEEdges(componentcveedges interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComponentCVEEdges", reflect.TypeOf((*MockIndexer)(nil).AddComponentCVEEdges), componentcveedges)
}

// Count mocks base method.
func (m *MockIndexer) Count(q *v1.Query, opts ...blevesearch.SearchOption) (int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{q}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Count", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIndexerMockRecorder) Count(q interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{q}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIndexer)(nil).Count), varargs...)
}

// DeleteComponentCVEEdge mocks base method.
func (m *MockIndexer) DeleteComponentCVEEdge(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComponentCVEEdge", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComponentCVEEdge indicates an expected call of DeleteComponentCVEEdge.
func (mr *MockIndexerMockRecorder) DeleteComponentCVEEdge(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComponentCVEEdge", reflect.TypeOf((*MockIndexer)(nil).DeleteComponentCVEEdge), id)
}

// DeleteComponentCVEEdges mocks base method.
func (m *MockIndexer) DeleteComponentCVEEdges(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComponentCVEEdges", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComponentCVEEdges indicates an expected call of DeleteComponentCVEEdges.
func (mr *MockIndexerMockRecorder) DeleteComponentCVEEdges(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComponentCVEEdges", reflect.TypeOf((*MockIndexer)(nil).DeleteComponentCVEEdges), ids)
}

// MarkInitialIndexingComplete mocks base method.
func (m *MockIndexer) MarkInitialIndexingComplete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkInitialIndexingComplete")
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkInitialIndexingComplete indicates an expected call of MarkInitialIndexingComplete.
func (mr *MockIndexerMockRecorder) MarkInitialIndexingComplete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkInitialIndexingComplete", reflect.TypeOf((*MockIndexer)(nil).MarkInitialIndexingComplete))
}

// NeedsInitialIndexing mocks base method.
func (m *MockIndexer) NeedsInitialIndexing() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NeedsInitialIndexing")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NeedsInitialIndexing indicates an expected call of NeedsInitialIndexing.
func (mr *MockIndexerMockRecorder) NeedsInitialIndexing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NeedsInitialIndexing", reflect.TypeOf((*MockIndexer)(nil).NeedsInitialIndexing))
}

// Search mocks base method.
func (m *MockIndexer) Search(q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{q}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Search", varargs...)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockIndexerMockRecorder) Search(q interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{q}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIndexer)(nil).Search), varargs...)
}