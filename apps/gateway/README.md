# Go Gateway

Go HTTP API Gateway for E-director.

Responsibilities:

- Expose browser-facing HTTP APIs.
- Accept JSON requests from the Vue frontend.
- Orchestrate auth, job, and generation workflows.
- Manage job/session state through MySQL and Redis.
- Keep Python AI Worker hidden from the browser.
