# Contributing to psirent

This document outlines our practices in the codebase.

## Folder Structure

### **`peer/send/` Folder**

- Use the `peer/send/` folder for implementing **peer-related functions**.
- These functions define operations from the **peer's perspective** and handle interactions initiated by peers.

### **`coordinator/receive/` Folder**

- Use the `coordinator/receive/` folder for implementing **coordinator-related functions**.
- These functions are answers to those in `peer/send/` folder from the **coordinator's perspective**.

## Naming Conventions

### **Importing `internal/errors`**

- When working with the `internal/errors` package, alias it as `errors2` to avoid conflicts with Go's standard `errors`
  package.
  ```go
  import (
      "errors"
      errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
  )
  // ...

### `io.Readers` and `io.Writers`

- Use *p* prefix before any of the `io` interfaces to denote a peer connection from the **coordinator's perspective**
  ```go
  func coordinatorFunction(pr io.Reader) {
      buf := make([]byte, 1024)
      _, _ = pr.Read(buf)
      fmt.Printf("message from peer: %v\n", string(buf))
      // ... 
  }
- Similarly, Use *c* prefix to denote a point of connection to the coordinator from the **peer's perspective**
  ```go
  func peerFunction(cw io.Writer) {
      _, _ = cw.Write([]byte("Hello, to coordinator!"))
      // ... 
  }

## Communication Logic

### internal/coms

Use the internal/coms package to share utility functions and abstractions between the **coordinator** and **peers**.

### Buffers

Always read/write strings as if they were delimited with a newline. Do **NOT** read arbitrary amount of bytes from a
buffer.