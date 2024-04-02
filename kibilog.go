package gokibilog

import (
	"fmt"
	"sync"
)

var once sync.Once
var instance *Kibilog

// Kibilog is singleton entity. Use [GetInstance] to get this.
type Kibilog struct {
	mu        sync.Mutex
	authToken string
	pools     map[string]*LogPool
}

// SetAuthToken registers the user's api token required to send messages to Kibilog.com
func (k *Kibilog) SetAuthToken(authToken string) {
	k.authToken = authToken
}

func (k *Kibilog) getAuthToken() string {
	return k.authToken
}

// AddLogPool allows you to register another [LogPool].
func (k *Kibilog) AddLogPool(pool *LogPool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.pools[pool.getLogId()] = pool
}

// GetLogPoolById returns [LogPool] by its LogID if it was previously set.
//
// Otherwise, it returns an error.
func (k Kibilog) GetLogPoolById(logId string) (logPool *LogPool, err error) {
	logPool, ok := k.pools[logId]
	if !ok {
		return nil, fmt.Errorf("LogPool with id \"%s\" is not found!", logId)
	}
	return logPool, nil
}

// SendMessages sends all messages that were previously posted in all registered [LogPool].
func (k *Kibilog) SendMessages() (errs []error) {
	for _, pool := range k.pools {
		poolErrs := validatePool(pool)
		if len(poolErrs) > 0 {
			errs = append(errs, poolErrs...)
		}
		// TODO send pool
	}
	return errs
}

// GetInstance allows you to get a single instance of [Kibilog].
func GetInstance() *Kibilog {
	once.Do(func() {
		instance = new(Kibilog)
		instance.pools = make(map[string]*LogPool)
	})
	return instance
}
