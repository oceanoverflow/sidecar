package loadbalancing

import (
	"sync"
)

type weighted struct {
	Server          string
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

func (w *weighted) fail() {
	w.EffectiveWeight -= w.Weight
	if w.EffectiveWeight < 0 {
		w.EffectiveWeight = 0
	}
}

func nextWeighted(servers []*weighted) (best *weighted) {
	total := 0

	for i := 0; i < len(servers); i++ {
		w := servers[i]

		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return
}

// WeightedServers is a collection of weighted server
type WeightedServers struct {
	sync.RWMutex
	servers    []*weighted
	serversMap map[string]bool
	n          int
}

func New() *WeightedServers {
	return &WeightedServers{
		serversMap: make(map[string]bool),
		n:          0,
	}
}

// Add add one server to the collection
func (ws *WeightedServers) Add(server string, weight int) {
	ws.Lock()
	defer ws.Unlock()
	if _, ok := ws.serversMap[server]; !ok {
		w := &weighted{
			Server:          server,
			Weight:          weight,
			EffectiveWeight: weight,
		}
		if ws.n >= len(ws.servers) {
			ws.servers = append(ws.servers, w)
		} else {
			ws.servers[ws.n] = w
		}
		ws.serversMap[server] = true
		ws.n++
	}
}

// Next get the next server for load balancing
func (ws *WeightedServers) Next() string {
	i := ws.nextWeighted()
	if i == nil {
		return ""
	}
	return i.Server
}

func (ws *WeightedServers) Len() int {
	ws.RLock()
	defer ws.RUnlock()
	return ws.n
}

func (ws *WeightedServers) nextWeighted() *weighted {
	if ws.n == 0 {
		return nil
	}
	if ws.n == 1 {
		return ws.servers[0]
	}
	return nextWeighted(ws.servers)
}
