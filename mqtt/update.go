package mqtt

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gomqtt/client"
	"github.com/gomqtt/packet"
)

// UpdateStatus is emitted by updateOne and Update.
type UpdateStatus struct {
	Progress float64
	Error    error
}

// Update will concurrently perform a firmware update and block until all devices
// have updated or returned errors. If a callback is provided it will be called
// with the current status of the update.
func Update(url string, baseTopics []string, firmware []byte, timeout time.Duration, callback func(string, *UpdateStatus)) error {
	// prepare table
	table := make(map[string]*UpdateStatus)

	// prepare wait group
	var wg = sync.WaitGroup{}

	// callback mutex
	var mutex sync.Mutex

	for _, baseTopic := range baseTopics {
		// create status
		table[baseTopic] = &UpdateStatus{}

		// add process
		wg.Add(1)

		// run a single update in background
		go func(bt string) {
			// begin update
			err := updateOne(url, bt, firmware, timeout, func(progress float64) {
				// lock mutex
				mutex.Lock()
				defer mutex.Unlock()

				// update progress
				table[bt].Progress = progress

				// call callback if provided
				if callback != nil {
					callback(bt, table[bt])
				}
			})
			if err != nil {
				// lock mutex
				mutex.Lock()
				defer mutex.Unlock()

				// update error
				table[bt].Error = err

				// call callback if provided
				if callback != nil {
					callback(bt, table[bt])
				}
			}

			// remove from wait group
			wg.Done()
		}(baseTopic)
	}

	// wait for all updates to complete
	wg.Wait()

	// return first error
	for _, us := range table {
		if us.Error != nil {
			return us.Error
		}
	}

	return nil
}

func updateOne(url, baseTopic string, firmware []byte, timeout time.Duration, progress func(float64)) error {
	// prepare channels
	requests := make(chan int)
	errs := make(chan error)

	// create client
	cl := client.New()

	// set callback
	cl.Callback = func(msg *packet.Message, err error) {
		// send errors
		if err != nil {
			errs <- err
			return
		}

		// otherwise convert the chunk request
		n, err := strconv.ParseInt(string(msg.Payload), 10, 0)
		if err != nil {
			errs <- err
			return
		}

		// check size
		if n <= 0 {
			errs <- fmt.Errorf("invalid chunk request of size %d", n)
			return
		}

		// send chunk request
		requests <- int(n)
	}

	// connect to the broker using the provided url
	cf, err := cl.Connect(client.NewConfig(url))
	if err != nil {
		return err
	}

	// wait for ack
	err = cf.Wait(timeout)
	if err != nil {
		return err
	}

	// make sure client gets closed
	defer cl.Close()

	// subscribe to next chunk request topic
	sf, err := cl.Subscribe(baseTopic+"/naos/update/request", 0)
	if err != nil {
		return err
	}

	// wait for ack
	err = sf.Wait(timeout)
	if err != nil {
		return err
	}

	// begin update process by sending the size of the firmware
	pf, err := cl.Publish(baseTopic+"/naos/update/begin", []byte(strconv.Itoa(len(firmware))), 0, false)
	if err != nil {
		return err
	}

	// wait for ack
	err = pf.Wait(timeout)
	if err != nil {
		return err
	}

	// prepare counters
	total := 0

	for {
		// the max size of the next requested chunk
		maxSize := 0

		// wait for error or request
		select {
		case err := <-errs:
			return err
		case <-time.After(timeout):
			return fmt.Errorf("update request timeout")
		case maxSize = <-requests:
			// continue
		}

		// calculate remaining bytes
		remaining := len(firmware) - total

		// check if done
		if remaining == 0 {
			// send finish
			pf, err := cl.Publish(baseTopic+"/naos/update/finish", nil, 0, false)
			if err != nil {
				return err
			}

			// wait for ack
			err = pf.Wait(timeout)
			if err != nil {
				return err
			}

			// disconnect client
			cl.Disconnect()

			return nil
		}

		// check chunk size
		if maxSize > remaining {
			// prevent overflow
			maxSize = remaining
		} else if maxSize > 4096 {
			// limit to 4096 bytes
			maxSize = 4096
		}

		// write chunk
		pf, err := cl.Publish(baseTopic+"/naos/update/write", firmware[total:total+maxSize], 0, false)
		if err != nil {
			return err
		}

		// wait for ack
		err = pf.Wait(timeout)
		if err != nil {
			return err
		}

		// adjust counter
		total += maxSize

		// update progress if available
		if progress != nil {
			progress(1.0 / float64(len(firmware)) * float64(total))
		}
	}
}