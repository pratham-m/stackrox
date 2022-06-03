// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"
	"time"

	"github.com/lib/pq"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTablePoliciesStmt holds the create statement for table `policies`.
	CreateTablePoliciesStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists policies (
                   Id varchar,
                   Name varchar UNIQUE,
                   Description varchar,
                   Disabled bool,
                   Categories text[],
                   LifecycleStages int[],
                   Severity integer,
                   EnforcementActions int[],
                   LastUpdated timestamp,
                   SORTName varchar,
                   SORTLifecycleStage varchar,
                   SORTEnforcement bool,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		GormModel: (*Policies)(nil),
		Indexes:   []string{},
		Children:  []*postgres.CreateStmts{},
	}

	// PoliciesSchema is the go schema for table `policies`.
	PoliciesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("policies")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.Policy)(nil)), "policies")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_POLICIES, "policy", (*storage.Policy)(nil)))
		RegisterTable(schema, CreateTablePoliciesStmt)
		return schema
	}()
)

const (
	PoliciesTableName = "policies"
)

// Policies holds the Gorm model for Postgres table `policies`.
type Policies struct {
	Id                 string           `gorm:"column:id;type:varchar;primaryKey"`
	Name               string           `gorm:"column:name;type:varchar;unique"`
	Description        string           `gorm:"column:description;type:varchar"`
	Disabled           bool             `gorm:"column:disabled;type:bool"`
	Categories         *pq.StringArray  `gorm:"column:categories;type:text[]"`
	LifecycleStages    *pq.Int32Array   `gorm:"column:lifecyclestages;type:int[]"`
	Severity           storage.Severity `gorm:"column:severity;type:integer"`
	EnforcementActions *pq.Int32Array   `gorm:"column:enforcementactions;type:int[]"`
	LastUpdated        *time.Time       `gorm:"column:lastupdated;type:timestamp"`
	SORTName           string           `gorm:"column:sortname;type:varchar"`
	SORTLifecycleStage string           `gorm:"column:sortlifecyclestage;type:varchar"`
	SORTEnforcement    bool             `gorm:"column:sortenforcement;type:bool"`
	Serialized         []byte           `gorm:"column:serialized;type:bytea"`
}