package tfdyn

import (
	"context"
	"testing"

	"github.com/databricks/cli/bundle/config/resources"
	"github.com/databricks/cli/bundle/internal/tf/schema"
	"github.com/databricks/cli/libs/dyn"
	"github.com/databricks/cli/libs/dyn/convert"
	"github.com/databricks/databricks-sdk-go/service/catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertTable(t *testing.T) {
	src := resources.Table{
		Name:        "test_table",
		CatalogName: "main",
		SchemaName:  "default",
		TableType:   catalog.TableTypeManaged,
		DataSourceFormat: catalog.DataSourceFormatDelta,
		Comment:     "Test table",
		Columns: []catalog.ColumnInfo{
			{
				Name:     "id",
				TypeName: catalog.ColumnTypeNameLong,
				TypeText: "BIGINT",
				Position: 0,
				Nullable: false,
			},
			{
				Name:     "name",
				TypeName: catalog.ColumnTypeNameString,
				TypeText: "STRING",
				Position: 1,
				Nullable: true,
			},
		},
		Grants: []resources.Grant{
			{
				Privileges: []string{"SELECT"},
				Principal:  "test-group",
			},
		},
	}

	vin, err := convert.FromTyped(src, dyn.NilValue)
	require.NoError(t, err)

	ctx := context.Background()
	out := schema.NewResources()
	err = tableConverter{}.Convert(ctx, "my_table", vin, out)
	require.NoError(t, err)

	// Assert that the table resource was created.
	assert.Equal(t, map[string]any{
		"catalog_name":        "main",
		"schema_name":         "default",
		"name":                "test_table",
		"table_type":          "MANAGED",
		"data_source_format":  "DELTA",
		"comment":             "Test table",
	}, out.Table["my_table"])

	// Assert that the grants resource was created.
	assert.Equal(t, &schema.ResourceGrants{
		Table: "${databricks_table.my_table.id}",
		Grant: []schema.ResourceGrantsGrant{
			{
				Privileges: []string{"SELECT"},
				Principal:  "test-group",
			},
		},
	}, out.Grants["table_my_table"])
}
