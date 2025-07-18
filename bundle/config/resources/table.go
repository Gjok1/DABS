package resources

import (
	"context"
	"net/url"
	"strings"

	"github.com/databricks/databricks-sdk-go/apierr"

	"github.com/databricks/cli/libs/log"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/marshal"
	"github.com/databricks/databricks-sdk-go/service/catalog"
)

// TableColumn represents a table column for Terraform
type TableColumn struct {
	Name             string `json:"name"`
	Type             string `json:"type,omitempty"`             // For SQL tables
	TypeName         string `json:"type_name,omitempty"`        // For regular tables
	TypeText         string `json:"type_text,omitempty"`        // For regular tables
	TypeJson         string `json:"type_json"`
	Position         int    `json:"position,omitempty"`         // For regular tables
	Nullable         bool   `json:"nullable,omitempty"`
	Comment          string `json:"comment,omitempty"`
	PartitionIndex   int    `json:"partition_index,omitempty"`
	TypePrecision    int    `json:"type_precision,omitempty"`
	TypeScale        int    `json:"type_scale,omitempty"`
	TypeIntervalType string `json:"type_interval_type,omitempty"`
}

type Table struct {
	// List of grants to apply on this table.
	Grants []Grant `json:"grants,omitempty"`

	// Full name of the table (catalog_name.schema_name.table_name). This value is read from
	// the terraform state after deployment succeeds.
	ID string `json:"id,omitempty" bundle:"readonly"`

	// Table name
	Name string `json:"name"`

	// Parent catalog for table
	CatalogName string `json:"catalog_name"`

	// Parent schema relative to its parent catalog
	SchemaName string `json:"schema_name"`

	// Type of table
	TableType catalog.TableType `json:"table_type,omitempty"`

	// Data source format
	DataSourceFormat catalog.DataSourceFormat `json:"data_source_format,omitempty"`

	// Columns for the table
	Columns []TableColumn `json:"column,omitempty"`

	// User-provided free-form text description
	Comment string `json:"comment,omitempty"`

	// A map of key-value properties attached to the securable
	Properties map[string]string `json:"properties,omitempty"`

	// Storage Location URL (full path) for the table
	StorageLocation string `json:"storage_location,omitempty"`

	// SQL warehouse ID for table operations
	WarehouseId string `json:"warehouse_id,omitempty"`

	ModifiedStatus ModifiedStatus `json:"modified_status,omitempty" bundle:"internal"`
	URL            string         `json:"url,omitempty" bundle:"internal"`
}

func (t *Table) Exists(ctx context.Context, w *databricks.WorkspaceClient, fullName string) (bool, error) {
	log.Tracef(ctx, "Checking if table with fullName=%s exists", fullName)

	_, err := w.Tables.Get(ctx, catalog.GetTableRequest{
		FullName: fullName,
	})
	if err != nil {
		log.Debugf(ctx, "table with full name %s does not exist: %v", fullName, err)

		if apierr.IsMissing(err) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func (*Table) ResourceDescription() ResourceDescription {
	return ResourceDescription{
		SingularName:  "table",
		PluralName:    "tables",
		SingularTitle: "Table",
		PluralTitle:   "Tables",
	}
}

func (t *Table) InitializeURL(baseURL url.URL) {
	if t.ID == "" {
		return
	}
	baseURL.Path = "explore/data/" + strings.ReplaceAll(t.ID, ".", "/")
	t.URL = baseURL.String()
}

func (t *Table) GetURL() string {
	return t.URL
}

func (t *Table) GetName() string {
	return t.Name
}

func (t *Table) UnmarshalJSON(b []byte) error {
	return marshal.Unmarshal(b, t)
}

func (t Table) MarshalJSON() ([]byte, error) {
	return marshal.Marshal(t)
}
