# quasar-mcp-server

 * quasar MCP Server

```shell
$ go install https://github.com/itsubaki/quasar-mcp-server@latest
```

```shell
$ gcloud run services describe ${SERVICE_NAME} --project ${PROJECT_ID} --format 'value(status.url)'
$ gcloud auth print-identity-token
```

```json
{
  "mcp": {
    "servers": {
      "quasar": {
        "command": "quasar-mcp-server",
        "env": {
          "BASE_URL": "YOUR_CLOUD_RUN_URL",
          "IDENTITY_TOKEN": "YOUR_IDENTITY_TOKEN"
        }
      }
    }
  }
}
```
