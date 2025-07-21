package tfdyn

import (
	"context"
	"fmt"

	"github.com/databricks/cli/bundle/internal/tf/schema"
	"github.com/databricks/cli/libs/dyn"
	"github.com/databricks/cli/libs/dyn/convert"
	"github.com/databricks/cli/libs/log"
)

func convertTableResource(ctx context.Context, vin dyn.Value) (dyn.Value, error) {
	// Normalize the output value to the target schema.
	v, diags := convert.Normalize(schema.ResourceTable{}, vin)
	for _, diag := range diags {
		log.Debugf(ctx, "table normalization diagnostic: %s", diag.Summary)
	}

	return v, nil
}

func convertSqlTableResource(ctx context.Context, vin dyn.Value) (dyn.Value, error) {
	// Normalize the output value to the SQL table schema.
	v, diags := convert.Normalize(schema.ResourceSqlTable{}, vin)
	for _, diag := range diags {
		log.Debugf(ctx, "sql table normalization diagnostic: %s", diag.Summary)
	}

	return v, nil
}

type tableConverter struct{}

func (tableConverter) Convert(ctx context.Context, key string, vin dyn.Value, out *schema.Resources) error {
	// Check if warehouse_id is specified to determine if we should use SQL table
	warehouseId, err := dyn.GetByPath(vin, dyn.NewPath(dyn.Key("warehouse_id")))
	useSqlTable := err == nil && warehouseId.IsValid() && warehouseId.MustString() != ""

	if useSqlTable {
		// Use SQL table for warehouse-enabled operations
		vout, err := convertSqlTableResource(ctx, vin)
		if err != nil {
			return err
		}

		// Add the converted resource to the SQL table output.
		out.SqlTable[key] = vout.AsAny()

		// Configure grants for this SQL table resource.
		if grants := convertGrantsResource(ctx, vin); grants != nil {
			grants.Table = fmt.Sprintf("${databricks_sql_table.%s.id}", key)
			out.Grants["sql_table_"+key] = grants
		}
	} else {
		// Use regular table
		vout, err := convertTableResource(ctx, vin)
		if err != nil {
			return err
		}

		// Add the converted resource to the regular table output.
		out.Table[key] = vout.AsAny()

		// Configure grants for this resource.
		if grants := convertGrantsResource(ctx, vin); grants != nil {
			grants.Table = fmt.Sprintf("${databricks_table.%s.id}", key)
			out.Grants["table_"+key] = grants
		}
	}

	return nil
}

func init() {
	registerConverter("tables", tableConverter{})
}
