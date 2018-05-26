# ToDo
Things needs to be done

## Dependencies
* github.com/gorilla/mux
    * HTTP router
* github.com/google/uuid
    * UUID
* github.com/peterbourgon/diskv
    * Disk storage
* github.com/pkg/errors
    * Errors with stack
* github.com/stretchr/testify/assert
    * Unit test assertion

## Project structures
### data
It's a folder to store the data.
Each filename is the identifier of ToDo document

### handlers
It's a package which includes http handler to manage the following resources:
``` go
POST    /v1/todo
POST    /v1/todo/{id}/tasks
PUT	/v1/todo/{id}/task/{taskID}/complete
GET	/v1/todo?search={search}&skip={skip}&limit={limit}
GET	/v1/todo/{id}
PUT	/v1/todo/{id}
DELETE  /v1/todo/{id}
DELETE	/v1/todo/{id}/task{taskID}
```
It includes unit test where it utilizes mock-up repository

### models
It's a package for request and response payload

### repositories
It's a package for repository (data access).  It contains an interface, file storage implementation, 
and mock-up repository (used for unit test).
File storage implements *diskv* where each file
contains each todo record, which includes list of tasks



