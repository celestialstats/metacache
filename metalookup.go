package metacache

import ()

type MetaLookup struct {
	Parameters map[string]interface{}
	Function   func(map[string]interface{}) map[string]string
}
