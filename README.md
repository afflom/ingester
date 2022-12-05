# Emporous Ingester

## Summary
Ingest content and automatically label it with a normalized schema. 

## Approach

Ingester performs a two part process for labeling data. 

1. Ingester references a schema and maps keys/values from within content to the referenced schema.

2. Once the labels are constructed, they are then converted to the Emporous normalized content schema. 

Ingester should maintain a list of supported filetypes that it can parse for key/values.

