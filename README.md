# go timeout examples

Just exploring some examples of how to do timeouts in go.

Run the server code:

```
go run server.go
```

Try out the various callers to verify they error w/ timeout after 2s:

```
time go run caller1-time-after.go 
```

```
time go run caller2-context-with-timeout.go 

```
time go run caller3-waitgroup.go 

```
time go run caller4-http-req-with-context.go
```
