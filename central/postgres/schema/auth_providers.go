// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableAuthProvidersStmt holds the create statement for table `auth_providers`.
	CreateTableAuthProvidersStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists auth_providers (
                   Id varchar,
                   Name varchar UNIQUE,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}

	// AuthProvidersSchema is the go schema for table `auth_providers`.
	AuthProvidersSchema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("auth_providers")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.AuthProvider)(nil)), "auth_providers")
		globaldb.RegisterTable(schema)
		return schema
	}()
)