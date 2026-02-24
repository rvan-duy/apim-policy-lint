# apim-policy-lint

Sample APIM policy fixtures are available in:

- `testdata/policies/basic-pass-through.xml`
- `testdata/policies/add-correlation-id.xml`
- `testdata/policies/rate-limit-by-subscription.xml`

Example usage:

```bash
go run . validate testdata/policies/basic-pass-through.xml
```

`validate` now parses the XML policy and emits a raw JSON document:

- `schemaVersion`: output contract version (`v1`)
- `sourceFile`: input XML path
- `raw`: generic XML AST (name/attrs/text/children)

Write JSON to a file while still printing to stdout:

```bash
go run . validate testdata/policies/basic-pass-through.xml --out /tmp/policy.json
```
