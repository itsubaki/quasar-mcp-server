# quasar-mcp-server

 * quasar MCP Server  
 * Run code written in OpenQASM format 

## Running a Bell State Circuit with OpenQASM

![claude desktop](claude_desktop.png)

## Installation and Environments

```shell
go install github.com/itsubaki/quasar-mcp-server@latest
gcloud run services describe ${SERVICE_NAME} --project ${PROJECT_ID} --format 'value(status.url)'
gcloud auth print-identity-token
```

