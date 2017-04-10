package metacache

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"strconv"
	"time"
)

type MetaCache struct {
	data          map[string]map[string]string
	updateChannel chan map[string]MetaLookup
	ValidMinutes  int
}

func NewMetaCache(ValidMinutes, MaxQueue int) *MetaCache {
	mc := new(MetaCache)
	mc.ValidMinutes = ValidMinutes
	mc.data = make(map[string]map[string]string)
	mc.updateChannel = make(chan map[string]MetaLookup, MaxQueue)
	go mc.Update()
	return mc
}

func (cache *MetaCache) CheckAndUpdate(key string, lookup MetaLookup) {
	if _, ok := cache.data[key]; ok {
		log.Println("Key", key, "exists. Checking if still valid...")
		if _, ok := cache.data[key]["Updated"]; ok {
			log.Println("\tHas Updated Attribute...")
			cacheTime, _ := strconv.ParseInt(cache.data[key]["Updated"], 10, 64)
			curTime := time.Now().UTC().Unix()
			olderInvalid := curTime - int64(cache.ValidMinutes*60)
			if olderInvalid > cacheTime {
				log.Println("\t\tOut of Date...")
			} else {
				log.Println("\t\tUp To Date (", (curTime - cacheTime), "s old)...")
				// We're up to date so return early...
				return
			}
		} else {
			log.Println("\tMissing Updated Attribute...")
		}
	} else {
		log.Println("Key", key, "does not exist.")
	}
	log.Println("Queueing to update Key:", key)
	// Everything else dropped through, so update
	cache.updateChannel <- map[string]MetaLookup{
		key: lookup,
	}
}

func (cache *MetaCache) PrintData() {
	spew.Dump(cache.data)
}

func (cache *MetaCache) Update() {
	for theMap := range cache.updateChannel {
		spew.Dump(theMap)
		for key, lookup := range theMap {
			cache.data[key] = lookup.Function(lookup.Parameters)
			cache.data[key]["Updated"] = strconv.FormatInt(time.Now().UTC().Unix(), 10)
		}
	}
}
