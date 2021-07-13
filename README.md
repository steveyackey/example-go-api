# example-go-api

This repository shows a very basic example of an API written in Go. There are two separate versions:
  - Standard library only
  - Standard libraray + httprouter

If this were a larger, production API, I'd probably separate out the handlers to their own file(s), routers to their own, logging or middleware in general to its own (and potentially either write a logging wrapper, or use something like zap), and split things. Though, I find it helpful to start simply, and refactor as it grows.