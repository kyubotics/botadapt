package adapters

import (
	"log"
	"strings"

	"github.com/richardchien/botadapt/adapters/coolqhttpapi"
	"github.com/richardchien/botadapt/t"
)

type AdapterMap map[string]t.Adapter

var (
	_adapters = AdapterMap{}
)

func init() {
	_adapters[coolqhttpapi.VIA] = coolqhttpapi.Adapter{}

	vias := []string{}
	for k, _ := range _adapters {
		vias = append(vias, k)
	}
	log.Println("Loaded adapters: " + strings.Join(vias, ", "))
}

func Find(event *t.Event) t.Adapter {
	via, ok := event.MessageSource["via"]
	if ok {
		a := _adapters[via]
		if a != nil {
			log.Println("Found adapter for via: " + via)
		}
		return a
	}
	return nil
}
