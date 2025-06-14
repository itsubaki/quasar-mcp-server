# quasar-mcp-server

 * quasar MCP Server  
 * Run code written in OpenQASM format 

## Running a Bell State Circuit with OpenQASM

![claude desktop](claude_desktop.png)

## Factoring with Shor's Algorithm

![claude desktop](claude_desktop_shor.png)

## Installation and Environments

 1. Deploy [quasar](https://github.com/itsubaki/quasar) to Cloud Run.
 1. Deploy quasar-mcp-server to Cloud Run.
 1. Configure `settings.json`.

```shell
$ make build
$ make deploy
```

```json
{
    "mcp": {
        "servers": {
            "quasar": {
                "type": "http",
                "url": "https://${YOUR_CLOUD_RUN_SERVICE_URL}.run.app/mcp",
            }
        }
    }
}
```
