package voting

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Result struct {
	Shortcode string `json:"shortcode"`
	NumVotes  int    `json:"votes"`
}

type ByVotes []*Result

func (s ByVotes) Len() int      { return len(s) }
func (s ByVotes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByVotes) Less(i, j int) bool {
	return s[i].NumVotes > s[j].NumVotes
}

type Poll interface {
	Vote(choice string) error
	Results() ([]*Result, error)
}

type onFilePoll struct {
	sync.RWMutex
	file    *os.File
	counter *prometheus.CounterVec
}

func (p *onFilePoll) Vote(choice string) error {
	votes, err := p.readMap()
	if err != nil {
		return err
	}
	votes[choice]++
	err = p.writeMap(votes)
	if err != nil {
		return err
	}
	p.counter.With(prometheus.Labels{"emoji": choice}).Inc()
	log.Printf("Voted for [%s], which now has a total of [%d] votes", choice, votes[choice])
	return nil
}

func (p *onFilePoll) Results() ([]*Result, error) {
	votes, err := p.readMap()
	if err != nil {
		return nil, err
	}
	results := make([]*Result, len(votes))
	i := 0
	for emoji, numVotes := range votes {
		results[i] = &Result{emoji, numVotes}
		i++
	}
	sort.Sort(ByVotes(results))
	return results, nil
}

func (p *onFilePoll) readMap() (votes map[string]int, err error) {
	var vj []byte
	p.RLock()
	_, err = p.file.Seek(0, io.SeekStart)
	if err == nil {
		vj, err = io.ReadAll(p.file)
	}
	p.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("failed to read votes from %q: %w", p.file.Name(), err)
	}
	if len(vj) > 0 {
		err = json.Unmarshal(vj, &votes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal votes from %q: %w", p.file.Name(), err)
		}
	}
	if votes == nil {
		votes = make(map[string]int)
	}
	return votes, nil
}

func (p *onFilePoll) writeMap(votes map[string]int) error {
	vj, err := json.Marshal(votes)
	if err != nil {
		return fmt.Errorf("failed to marshal votes: %w", err)
	}
	p.Lock()
	_, err = p.file.Seek(0, io.SeekStart)
	if err == nil {
		_, err = p.file.Write(vj)
	}
	p.Unlock()
	if err != nil {
		err = fmt.Errorf("failed to write votes to %q: %w", p.file.Name(), err)
	}
	return err
}

type inMemoryPoll struct {
	votes map[string]int
	sync.RWMutex
	counter *prometheus.CounterVec
}

func (p *inMemoryPoll) Vote(choice string) error {
	p.Lock()
	defer p.Unlock()

	if p.votes[choice] > 0 {
		p.votes[choice] = p.votes[choice] + 1
	} else {
		p.votes[choice] = 1
	}
	p.counter.With(prometheus.Labels{"emoji": choice}).Inc()
	log.Printf("Voted for [%s], which now has a total of [%d] votes", choice, p.votes[choice])
	return nil
}

func (p *inMemoryPoll) Results() ([]*Result, error) {
	p.RLock()
	defer p.RUnlock()

	results := make([]*Result, len(p.votes))
	i := 0
	for emoji, numVotes := range p.votes {
		results[i] = &Result{emoji, numVotes}
		i++
	}
	sort.Sort(ByVotes(results))

	return results, nil
}

var counter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "emojivoto_votes_total",
	Help: "Number of emoji votes",
}, []string{"emoji"})

func NewInMemoryPoll() Poll {
	return &inMemoryPoll{
		votes:   make(map[string]int, 0),
		counter: counter,
	}
}

func NewOnFilePoll(path string) (Poll, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return nil, err
	}
	log.Printf("Storing votes in file %s", path)
	return &onFilePoll{
		file:    f,
		counter: counter,
	}, nil
}
