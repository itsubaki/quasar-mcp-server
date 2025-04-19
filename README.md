# quasar-mcp-server

 * quasar MCP Server

![claude desktop](claude_desktop.png)

## Install and environments

```shell
go install github.com/itsubaki/quasar-mcp-server@latest
gcloud run services describe ${SERVICE_NAME} --project ${PROJECT_ID} --format 'value(status.url)'
gcloud auth print-identity-token
```

