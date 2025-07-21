# Databricks Asset Bundles (DABS) - Table Deployment Guide

This repository contains the Databricks CLI with support for Databricks Asset Bundles (DABS), a powerful infrastructure-as-code solution for managing Databricks resources including tables, jobs, pipelines, and more.

## Table of Contents
- [Overview](#overview)
- [Building Locally](#building-locally)
- [Getting Started](#getting-started)
- [Table Configuration](#table-configuration)
- [Deployment](#deployment)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

Databricks Asset Bundles allow you to define and deploy tables as code using YAML configuration files. This approach provides:

- **Version Control**: Track changes to your table schemas and configurations
- **Environment Management**: Deploy the same table definitions across dev/staging/prod environments
- **Automated Deployment**: Use CI/CD pipelines to deploy table changes
- **Permissions Management**: Define grants and access controls declaratively

## Building Locally

### Prerequisites

- **Go 1.24.0 or later** (as specified in `cli/go.mod`)

### Building the CLI

1. **Navigate to the CLI directory**
   ```bash
   cd cli
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Build the CLI binary**
   ```bash
   go build
   ```

### Binary Location

After building, you'll find the binary at `./databricks` (in the cli directory).

### Cross-Platform Building

For specific platforms, you can use Go's cross-compilation:

```bash
# For Windows
GOOS=windows GOARCH=amd64 go build -o databricks.exe

# For macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o databricks-mac-intel

# For macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o databricks-mac-arm

# For Linux
GOOS=linux GOARCH=amd64 go build -o databricks-linux
```

### Quick Start

1. **Build the CLI**:
   ```bash
   cd cli
   go mod tidy
   go build
   ```

2. **Verify the build**:
   ```bash
   ./databricks --version
   ```

## Getting Started

### Prerequisites
- Databricks CLI installed and configured
- Access to a Databricks workspace with Unity Catalog enabled
- Appropriate permissions to create catalogs, schemas, and tables

### Initialize a Bundle

```bash
# Create a new bundle
databricks bundle init

# Navigate to your bundle directory
cd my-bundle
```

## Table Configuration

Tables are defined in your bundle's `databricks.yml` file under the `resources.tables` section. Here's a complete example:

```yaml
bundle:
  name: warehouse-table-example

resources:
  tables:
    customer_table:
      name: customer_table
      catalog_name: test
      schema_name: tabletest
      table_type: MANAGED
      data_source_format: DELTA
      comment: "Customer table with update support via warehouse"
      warehouse_id: "e88b8bde026c6b80"  # 🔑 KEY: This enables SQL warehouse operations
      column:
        - name: customer_id
          type: "BIGINT"
          type_json: '{"name":"customer_id","type":"long","nullable":false,"metadata":{}}'
          nullable: false
          position: 1
          comment: "Unique customer identifier"
        - name: first_name
          type: "STRING"
          type_json: '{"name":"first_name","type":"string","nullable":false,"metadata":{}}'
          nullable: false
          position: 2
          comment: "Customer first name"
        - name: last_name
          type: "STRING"
          type_json: '{"name":"last_name","type":"string","nullable":false,"metadata":{}}'
          nullable: false
          position: 3
          comment: "Customer last name"
        - name: email
          type: "STRING"
          type_json: '{"name":"email","type":"string","nullable":true,"metadata":{}}'
          nullable: true
          position: 4 
          comment: "Customer email address"
        - name: address
          type: "STRING"
          type_json: '{"name":"address","type":"string","nullable":true,"metadata":{}}'
          nullable: true
          position: 5
          comment: "Customer address"
        - name: system_id
          type: "STRING"
          type_json: '{"name":"system_id","type":"string","nullable":true,"metadata":{}}'
          nullable: true
          position: 6
          comment: "Id on Source System"
      properties:
        delta.autoOptimize.optimizeWrite: "true"
        delta.autoOptimize.autoCompact: "true"
```

### Table Types

#### Managed Tables
Managed tables store both data and metadata in the Unity Catalog:

```yaml
tables:
  my_managed_table:
    name: sales_data
    catalog_name: ${var.catalog_name}
    schema_name: ${var.schema_name}
    table_type: MANAGED
    data_source_format: DELTA
    # ... column definitions
```

#### External Tables
External tables reference data stored outside of Databricks:

```yaml
tables:
  my_external_table:
    name: external_data
    catalog_name: ${var.catalog_name}
    schema_name: ${var.schema_name}
    table_type: EXTERNAL
    data_source_format: PARQUET
    storage_location: "s3://my-bucket/external-data/"
    # ... column definitions
```

#### SQL Warehouse Tables
For SQL warehouse operations, include a `warehouse_id`:

```yaml
tables:
  customer_table:
    name: customer_table
    catalog_name: test
    schema_name: tabletest
    table_type: MANAGED
    data_source_format: DELTA
    comment: "Customer table with update support via warehouse"
    warehouse_id: "e88b8bde026c6b80"  # 🔑 KEY: This enables SQL warehouse operations
    column:
      - name: customer_id
        type: "BIGINT"
        type_json: '{"name":"customer_id","type":"long","nullable":false,"metadata":{}}'
        nullable: false
        position: 1
        comment: "Unique customer identifier"
      - name: first_name
        type: "STRING"
        type_json: '{"name":"first_name","type":"string","nullable":false,"metadata":{}}'
        nullable: false
        position: 2
        comment: "Customer first name"
      - name: last_name
        type: "STRING"
        type_json: '{"name":"last_name","type":"string","nullable":false,"metadata":{}}'
        nullable: false
        position: 3
        comment: "Customer last name"
      - name: email
        type: "STRING"
        type_json: '{"name":"email","type":"string","nullable":true,"metadata":{}}'
        nullable: true
        position: 4 
        comment: "Customer email address"
      - name: address
        type: "STRING"
        type_json: '{"name":"address","type":"string","nullable":true,"metadata":{}}'
        nullable: true
        position: 5
        comment: "Customer address"
      - name: system_id
        type: "STRING"
        type_json: '{"name":"system_id","type":"string","nullable":true,"metadata":{}}'
        nullable: true
        position: 6
        comment: "Id on Source System"
    properties:
      delta.autoOptimize.optimizeWrite: "true"
      delta.autoOptimize.autoCompact: "true"
```

### Column Definitions

Define table columns with detailed type information. **Important**: Always specify the `position` field to ensure proper column ordering.

```yaml
column:  # Note: use 'column' for SQL warehouse tables
  - name: customer_id
    type: "BIGINT"
    type_json: '{"name":"customer_id","type":"long","nullable":false,"metadata":{}}'
    nullable: false
    position: 1
    comment: "Unique customer identifier"
  - name: customer_name
    type: "STRING"
    type_json: '{"name":"customer_name","type":"string","nullable":false,"metadata":{}}'
    nullable: false
    position: 2
    comment: "Customer name"
```

**Note**: Column positions should start from 1 and be sequential. The position determines the column order in the table schema.

### Grants and Permissions

Define access controls directly in your table configuration:

```yaml
grants:
  - principal: "data-analysts"
    privileges: ["SELECT"]
  - principal: "data-engineers" 
    privileges: ["SELECT", "MODIFY"]
  - principal: "etl-service-principal"
    privileges: ["SELECT", "MODIFY", "CREATE"]
```

## Deployment

### Deploy Tables

Deploy your tables to the target environment:

```bash
# Deploy to default target
databricks bundle deploy

# Deploy with auto-approval (for CI/CD)
databricks bundle deploy --auto-approve

# Force deployment (override Git branch validation)
databricks bundle deploy --force

# Deploy with verbose output
databricks bundle deploy --verbose
```

### Environment-Specific Deployment

Configure different environments in your `databricks.yml`:

```yaml
targets:
  dev:
    mode: development
    default: true
    variables:
      catalog_name: dev_catalog
      schema_name: dev_analytics
      
  staging:
    mode: development
    variables:
      catalog_name: staging_catalog
      schema_name: staging_analytics
      
  prod:
    mode: production
    variables:
      catalog_name: prod_catalog
      schema_name: analytics
```

Deploy to specific environments:

```bash
# Deploy to development
databricks bundle deploy --target dev

# Deploy to staging
databricks bundle deploy --target staging

# Deploy to production
databricks bundle deploy --target prod
```

### Validation and Testing

Validate your bundle configuration before deployment:

```bash
# Validate bundle configuration
databricks bundle validate

# Run bundle tests
databricks bundle test
```

## Examples

### Complete Example: E-commerce Data Tables

```yaml
bundle:
  name: ecommerce-data

variables:
  catalog_name:
    default: ecommerce
  schema_name:
    default: raw_data

resources:
  schemas:
    raw_data_schema:
      name: ${var.schema_name}
      catalog_name: ${var.catalog_name}
      comment: "Raw e-commerce data schema"

  tables:
    # Customer dimension table
    customers:
      name: customers
      catalog_name: ${var.catalog_name}
      schema_name: ${var.schema_name}
      table_type: MANAGED
      data_source_format: DELTA
      comment: "Customer dimension table"
      columns:
        - name: customer_id
          type_name: LONG
          type_text: "BIGINT"
          nullable: false
          comment: "Unique customer identifier"
        - name: first_name
          type_name: STRING
          type_text: "STRING"
          nullable: false
        - name: last_name
          type_name: STRING
          type_text: "STRING"
          nullable: false
        - name: email
          type_name: STRING
          type_text: "STRING"
          nullable: true
        - name: created_at
          type_name: TIMESTAMP
          type_text: "TIMESTAMP"
          nullable: false
      properties:
        delta.autoOptimize.optimizeWrite: "true"
        delta.autoOptimize.autoCompact: "true"
      grants:
        - principal: "data-team"
          privileges: ["SELECT", "MODIFY"]

    # Orders fact table
    orders:
      name: orders
      catalog_name: ${var.catalog_name}
      schema_name: ${var.schema_name}
      table_type: MANAGED
      data_source_format: DELTA
      comment: "Orders fact table"
      columns:
        - name: order_id
          type_name: LONG
          type_text: "BIGINT"
          nullable: false
        - name: customer_id
          type_name: LONG
          type_text: "BIGINT"
          nullable: false
        - name: order_date
          type_name: DATE
          type_text: "DATE"
          nullable: false
        - name: total_amount
          type_name: DECIMAL
          type_precision: 10
          type_scale: 2
          type_text: "DECIMAL(10,2)"
          nullable: false
        - name: status
          type_name: STRING
          type_text: "STRING"
          nullable: false
      grants:
        - principal: "data-analysts"
          privileges: ["SELECT"]
        - principal: "data-engineers"
          privileges: ["SELECT", "MODIFY"]

targets:
  dev:
    mode: development
    default: true
    variables:
      catalog_name: dev_ecommerce
      
  prod:
    mode: production
    variables:
      catalog_name: prod_ecommerce
```

## Best Practices

### 1. Use Variables for Environment Configuration
- Define environment-specific values as variables
- Use different catalogs/schemas for dev/staging/prod
- Keep table structure consistent across environments

### 2. Schema Management
- Always create schemas before tables that depend on them
- Use descriptive comments for schemas and tables
- Follow consistent naming conventions

### 3. Column Design
- Define explicit column types and constraints
- Include meaningful comments for all columns
- Use appropriate data types for your use case

### 4. Access Control
- Define grants at the table level in your bundle
- Use service principals for automated processes
- Follow the principle of least privilege

### 5. Data Format Selection
- Use DELTA format for managed tables (default)
- Consider PARQUET for external tables
- Enable auto-optimization for Delta tables

### 6. Version Control
- Commit your bundle configuration to Git
- Use feature branches for table schema changes
- Tag releases for production deployments

### 7. CI/CD Integration
```bash
# Example CI/CD deployment script
#!/bin/bash
databricks bundle validate --target prod
databricks bundle deploy --target prod --auto-approve
```

### 8. Monitoring and Rollback
- Monitor deployment logs for errors
- Test table access after deployment
- Keep backup strategies for critical tables

## Troubleshooting

### Common Issues

1. **Permission Errors**: Ensure your service principal has appropriate Unity Catalog permissions
2. **Schema Dependencies**: Create schemas before tables that reference them
3. **Column Type Mismatches**: Verify column type definitions match your data
4. **Grant Failures**: Check that principals exist and have necessary permissions

### Getting Help

```bash
# View bundle help
databricks bundle --help

# View deploy command options
databricks bundle deploy --help

# Check bundle status
databricks bundle summary
```

For more information, consult the [Databricks documentation](https://docs.databricks.com/dev-tools/bundles/index.html) on Asset Bundles. 