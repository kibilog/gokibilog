package gokibilog

import (
	"encoding/json"
	"fmt"
)

func validatePool(pool *LogPool) (errs []error) {
	for k, m := range pool.messages {
		_, err := json.Marshal(m)
		if err != nil {
			errs = append(errs, fmt.Errorf("Error in the message for \"%s\": %s", pool.logId, err.Error()))
		}
		pool.messages[k] = nil
	}
	pool.removeNilMessages()

	return errs
}
