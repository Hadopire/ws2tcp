Really basic websocket to tcp proxy

```
Usage:
  -p string
    	specifies the port (default "8000")
```

On connection the client must provide the host he wishes to connect to.
```
{
  "host" : "hostname:port"
}
```
