# Lock manager

The lock manager provides a utility to interact with a lock service. 

## Lock service API

The manager expects the following REST API to be exposed by a lock manager:

### Create a lock

```http
POST /lock HTTP/1.1
Content-Type: application/json

{"workflow": "workflow-name"}
```

A success response will look like this:

```http
HTTP/1.1 200 Success

{
  "id": "e9205d4a-2545-4c58-86ae-d0e7b3ecec23",
  "status": "locked",
  "created": "2019-05-28T08:30:24.595Z",
  "lastChange": "2019-05-28T08:30:24.595Z"
}
```

### Get lock information

```http
GET /lock/?id=<lock-id> HTTP/1.1
```

A successful response will look like this:

```http
HTTP/1.1 200 Success

{
  "id": "e9205d4a-2545-4c58-86ae-d0e7b3ecec23",
  "status": "locked",
  "created": "2019-05-28T08:30:24.595Z",
  "lastChange": "2019-05-28T08:30:24.595Z"
}
```

### Release a lock

```http
PATCH /lock/?id=<lock-id> HTTP/1.1
```

A successful response will look like this:

```http
HTTP/1.1 200 Success

{
  "id": "e9205d4a-2545-4c58-86ae-d0e7b3ecec23",
  "status": "released",
  "created": "2019-05-28T08:30:24.595Z",
  "lastChange": "2019-05-28T08:30:24.595Z"
}
```

### Delete a lock

```http
DELETE /lock/?id=<lock-id> HTTP/1.1
```

A successful response will look like this:

```http
HTTP/1.1 200 Success

{}
```