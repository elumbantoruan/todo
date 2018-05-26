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

## Project structures
### data
It's a folder to store the data

### handlers
It's a package to manage resources (including unit test)

### models
It's a package for request and response payload

### repositories
It's a package for repository (data access).  It contains an interface, file storage implementation, 
and mockup repository (used for unit test)


