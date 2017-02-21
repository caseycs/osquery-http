# osquery-http
http endpoint for [osquery](https://github.com/facebook/osquery)

## Requirements

* go
* osquery

## Example

Quickstart:

```
go run main.go
curl 'http://localhost:8000/tables'
````

A bit more complex:

```
SECRET=foo go run main.go
curl 'http://localhost:8000/table/users?secret=foo'
````

And something else:

```
BIND=localhost:8000 SECRET=foo go run main.go
curl 'http://localhost:8000/query?q=select%20name,cmdline,total_size%20from%20processes%20order%20by%20total_size%20desc%20limit%2010&secret=foo'
````
