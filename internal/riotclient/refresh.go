package riotclient

import (
	"fmt"
	"sync"
	"time"
)

type refreshJob struct {
	jobID     string
	region    string
	puuid     string
	queuedAt  time.Time
	expiresAt time.Time
}

type refreshQueue struct {
	mu   sync.Mutex
	jobs map[string]*refreshJob // key: region:puuid
}

func newRefreshQueue() *refreshQueue {
	return &refreshQueue{jobs: make(map[string]*refreshJob)}
}

// RequestRefresh enqueues a summoner refresh with a 120s per-summoner
// cooldown. A repeat request inside the window returns the existing job and
// the remaining cooldown instead of queueing duplicates.
func (q *refreshQueue) RequestRefresh(region, puuid string) (jobID string, cooldownSec int, coalesced bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	key := region + ":" + puuid
	now := time.Now()
	if existing, ok := q.jobs[key]; ok && now.Before(existing.expiresAt) {
		return existing.jobID, int(existing.expiresAt.Sub(now).Seconds()), true
	}
	id := fmt.Sprintf("job-%s-%d", region, now.UnixNano())
	q.jobs[key] = &refreshJob{
		jobID:     id,
		region:    region,
		puuid:     puuid,
		queuedAt:  now,
		expiresAt: now.Add(120 * time.Second),
	}
	// Opportunistically evict expired entries so the map stays small.
	for k, job := range q.jobs {
		if now.After(job.expiresAt) {
			delete(q.jobs, k)
		}
	}
	return id, 120, false
}

// Pending counts jobs still inside their cooldown window.
func (q *refreshQueue) Pending() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	now := time.Now()
	n := 0
	for k, job := range q.jobs {
		if now.After(job.expiresAt) {
			delete(q.jobs, k)
			continue
		}
		n++
	}
	return n
}
