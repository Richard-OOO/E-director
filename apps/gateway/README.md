# Go Gateway

Go + Gin API Gateway for E-director.

Responsibilities:

- Expose browser-facing HTTP APIs.
- Convert AI Worker gRPC stream events to SSE events.
- Manage job/session state through Redis.
- Keep Python AI Worker hidden from the browser.
