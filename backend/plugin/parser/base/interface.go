package base

import (
	"fmt"
	"sync"

	"github.com/bytebase/bytebase/backend/plugin/db"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

var (
	mux                    sync.Mutex
	queryValidators        = make(map[storepb.Engine]ValidateSQLForEditorFunc)
	fieldMaskers           = make(map[storepb.Engine]GetMaskedFieldsFunc)
	changedResourcesGetter = make(map[storepb.Engine]ExtractChangedResourcesFunc)
	resourcesGetter        = make(map[storepb.Engine]ExtractResourceListFunc)
)

type ValidateSQLForEditorFunc func(string) bool
type GetMaskedFieldsFunc func(string, string, *db.SensitiveSchemaInfo) ([]db.SensitiveField, error)
type ExtractChangedResourcesFunc func(string, string, string) ([]SchemaResource, error)
type ExtractResourceListFunc func(string, string, string) ([]SchemaResource, error)

func RegisterQueryValidator(engine storepb.Engine, f ValidateSQLForEditorFunc) {
	mux.Lock()
	defer mux.Unlock()
	if _, dup := queryValidators[engine]; dup {
		panic(fmt.Sprintf("Register called twice %s", engine))
	}
	queryValidators[engine] = f
}

func RegisterGetMaskedFieldsFunc(engine storepb.Engine, f GetMaskedFieldsFunc) {
	mux.Lock()
	defer mux.Unlock()
	if _, dup := fieldMaskers[engine]; dup {
		panic(fmt.Sprintf("Register called twice %s", engine))
	}
	fieldMaskers[engine] = f
}

func RegisterExtractChangedResourcesFunc(engine storepb.Engine, f ExtractChangedResourcesFunc) {
	mux.Lock()
	defer mux.Unlock()
	if _, dup := changedResourcesGetter[engine]; dup {
		panic(fmt.Sprintf("Register called twice %s", engine))
	}
	changedResourcesGetter[engine] = f
}

func RegisterExtractResourceListFunc(engine storepb.Engine, f ExtractResourceListFunc) {
	mux.Lock()
	defer mux.Unlock()
	if _, dup := resourcesGetter[engine]; dup {
		panic(fmt.Sprintf("Register called twice %s", engine))
	}
	resourcesGetter[engine] = f
}

func ExtractSensitiveField(engine storepb.Engine, statement string, currentDatabase string, schemaInfo *db.SensitiveSchemaInfo) ([]db.SensitiveField, error) {
	if schemaInfo == nil {
		return nil, nil
	}

	f := fieldMaskers[engine]
	return f(statement, currentDatabase, schemaInfo)
}

// ValidateSQLForEditor validates the SQL statement for editor.
// We support the following SQLs:
// 1. EXPLAIN statement, except EXPLAIN ANALYZE
// 2. SELECT statement
// We also support CTE with SELECT statements, but not with DML statements.
func ValidateSQLForEditor(engine storepb.Engine, statement string) bool {
	f := queryValidators[engine]
	return f(statement)
}
