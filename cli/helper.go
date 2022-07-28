package cli

import (
	"gosail/logger"

	"github.com/desertbit/grumble"
)

var (
	log = logger.Logger()
)

func getValue(c *grumble.Context, name string, v interface{}) interface{} {
	var value interface{}
	switch v.(type) {
	case string:
		value = c.Flags.String(name)
		if value != v {
			return value
		} else {
			return c.Args.String(name)
		}
	case int:
		value = c.Flags.Int(name)
		if value != v {
			return value
		} else {
			return c.Args.Int(name)
		}
	case bool:
		value = c.Flags.Bool(name)
		if value != v {
			return value
		} else {
			return c.Args.Bool(name)
		}
	}
	return v
}
