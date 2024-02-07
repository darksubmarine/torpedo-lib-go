# Configuration

This package is useful to load configuration values from different data sources. The lib provides a simple key-value `Map`
with methods to fetch the loaded values. Also provides 2 loaders out of the box:

 - Yaml file loader
 - Azure AppConfig loader with Azure Key Vault integration

## How to use it?

This package let you set up the configuration data sources and once that has been loaded the values can be accessed 
via the `Map` interface. For instance:

```yaml
foo:
  bar:
    baz: "some config value"
```

```go
package main

import (
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/conf"
)

func main() {
	// 1. Read your config from different data sources
	cfg := conf.Load(false, conf.NewYamlLoader("config.yaml"))

	// 2. Fetch values from the Map
	val := cfg.FetchStringOrElse("defaultValue", "foo", "bar", "baz")

	fmt.Println(val)

}
```

