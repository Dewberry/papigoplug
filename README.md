# papigoplug

Go plugin utilities for the process-api. Equivalent to the Python implementation at https://github.com/Dewberry/papipyplug

## Usage

Any program leveraging papigoplug must initialize its log before invoking other papigoplug functions.
Otherwise the program will panic with a generic message when papigoplug tries to call log methods against a nil pointer.
`InitLog()` must be called once only. Calling it a second time will cause a `Log.Fatal()` indicating that the log
has already been initialized.

After initializing the log, it can be accessed like `papigoplug.Log.Info("hello, world!")`

```go
package main

import (
    papigoplug "github.com/Dewberry/papigoplug/papigoplug"
)

func main() {
    papigoplug.InitLog("info")
    // Now it's safe to call other papigoplug functions that might use Log.
}
```
