# quasar-mcp-server

An MCP server for running OpenQASM code

## Deployment and Configuration

 1. Deploy [quasar](https://github.com/itsubaki/quasar) to Cloud Run.
 1. Deploy **quasar-mcp-server** to Cloud Run.
 1. Configure the `settings.json`.

```shell
make build deploy
```

```json
{
    "mcp": {
        "servers": {
            "quasar": {
                "url": "https://${YOUR_CLOUD_RUN_SERVICE_URL}.run.app/mcp"
            }
        }
    }
}
```

## Invoking an Authenticated Cloud Run Service from localhost

```shell
make proxy
```

```json
{
    "mcp": {
        "servers": {
            "quasar": {
                "url": "http://127.0.0.1:3000/mcp"
            }
        }
    }
}
