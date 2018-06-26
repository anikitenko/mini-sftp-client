## Mini sFTP client API

### API was created using Swagger

In order to see detailed reference you should navigate 
to http://127.0.0.1:9000/api/v1

API has the following endpoints:
* PUT /connect
* DELETE /disconnect/{id}
* GET /download/{id}
* GET /getConnections
* GET /getLocalHomeDirectory
* GET /getLocalPathCompletion
* GET /getRemoteHomeDirectory/{id}
* GET /getRemotePathCompletion/{id}

If you are accessing from remote (not from localhost) then you need to have pin code. Pin code can be passed during each API request via header:
```
    curl -H "Pin-Code: 1234"
```
