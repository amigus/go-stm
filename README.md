# Simple Token Manager for Go

Simple Token Manager for Go provisions tokens then checks input strings against them.

Tokens are string values that can expire after _count_ uses and/or _timeout_ seconds.

It uses [Google uuid](https://github.com/google/uuid) to generate token values.

```golang
// The number of valid tokens that can exist at once
count := 2

// The number of times each token can be checked before it expires through use
uses := 3

// The amount of time each token can be checked before it expires through age
timeout := 2 * time.Seconds

tm := stm.UUIDTokenManager(count, uses, timeout)

t1 := tm.Get()
t2 := tm.Get()
t3 := tm.Get() // causes t1 to expire

if (tm.Check(t1)) { panic() } // t1 has expired

if (!tm.Check(t2)) { panic() } // t2 has 2 uses left
if (!tm.Check(t2)) { panic() } // t2 has 1 use left
if (!tm.Check(t2)) { panic() } // causes t2 to expire

if (tm.Check(t2)) { panic() } // t2 has expired

if (!tm.Check(t3)) { panic() } // t3 has 2 uses left

time.Sleep(timeout) // causes all tokens to expire

if (tm.Check(t3)) { panic() } // t3 has expired by age
```

## Usage

### Gin

The [gin](gin) directory contains functions that work with the [Gin Web Framework](https://gin-gonic.com/).

`TokenPublisher` exposes `Get()` as an endpoint.

`HeaderChecker` calls `Check()` on the value of a custom HTTP header and returns HTTP status 401 (Unauthorized) unless it returns `true`.

#### Examples

Check for a single, unlimited use, non-expiring token:

```golang
package main

import (
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/amigus/go-stm"
    stm_gin "github.com/amigus/go-stm/gin"
)

const (
    // HTTPHeaderName is the name of the header that contains the token
    HTTPHeaderName = "X-Token"
)

func main() {
    stm := stm.UUIDTokenManager(1, 0, 0) // One unlimited use non-expiring token

    fmt.Printf("The token is: %s\n", stm.Get())

    // Add the HeaderChecker middleware to the gin engine
    r := stm_gin.HeaderChecker(gin.Default(), stm, HTTPHeaderName)
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"output": "success"})
    })
    r.Run(":8080")
}
```

1. Run it
2. Get 401 Unauthorized accessing the resource without the token
3. Get 200 OK using the token

```bash
./main &
[1] 96580
The token is: 31ca3e1c-26ec-4368-9e91-4330cecfbce4
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (4 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080

curl localhost:8080
[GIN] 2025/01/07 - 11:42:34 | 401 |      66.774µs |             ::1 | GET      "/"
{"error":"Invalid token"}

curl -H 'X-Token: 31ca3e1c-26ec-4368-9e91-4330cecfbce4' localhost:8080
[GIN] 2025/01/07 - 11:42:06 | 200 |       77.61µs |             ::1 | GET      "/"
{"output":"success"}
```

Check for a single-use token from the `TokenPublisher` running on a UNIX socket:

```golang
package main

import (
    "github.com/gin-gonic/gin"

    "github.com/amigus/go-stm"
    stm_gin "github.com/amigus/go-stm/gin"
)

const (
    // HTTPHeaderName is the name of the header that contains the token
    HTTPHeaderName = "X-Token"
    // UnixSocketPath is the path to the unix socket
    UnixSocketPath = "stm.sock"
)

func main() {
    stm := stm.UUIDTokenManager(1, 1, 0) // One single-use non-expiring token

    // Start a TokenPublisher running on a unix socket
    go func() {
        stm_gin.TokenPublisher(gin.Default(), stm, "/").RunUnix(UnixSocketPath)
    }()

    // Add the HeaderChecker middleware to the gin engine
    r := stm_gin.HeaderChecker(gin.Default(), stm, HTTPHeaderName)
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"output": "success"})
    })
    r.Run(":8080")
    }
```

1. Run it
2. Get 401 Unauthorized accessing the resource without the token
3. Get the token and use it to get 200 OK

```bash
./main &
[1] 449
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func2 (4 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> github.com/amigus/stm/gin.TokenPublisher.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on unix:/stm.sock
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies for details.

curl localhost:8080
[GIN] 2025/01/07 - 11:46:52 | 401 |      64.534µs |             ::1 | GET      "/"
{"error":"Invalid token"}

curl -H "X-Token: $(curl --unix-socket ./stm.sock -s .)" localhost:8080
[GIN] 2025/01/07 - 11:48:23 | 200 |      45.557µs |                 | GET      "/"
[GIN] 2025/01/07 - 11:48:23 | 200 |     604.372µs |             ::1 | GET      "/"
{"output":"success"}
```
