package metacache

import ()

type MetaLookup struct {
	Parameters map[string]string
	Function   func(map[string]string) map[string]string
}
