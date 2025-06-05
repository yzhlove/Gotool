package act

import (
	"context"
	"fmt"
	"github.com/gofrs/flock"
	"github.com/yzhlove/Gotool/flocktest/module/log"
	"strconv"
	"time"
)

type actMode string

const (
	LockMode     actMode = "1"
	TryLockMode  actMode = "2"
	LockStatMode actMode = "3"
)

type Action interface {
	Act() error
}

func New(path, duration, mode string) Action {
	var act Action
	b := base{path: path, duration: duration}
	switch actMode(mode) {
	case LockMode:
		act = LockAct{b}
	case TryLockMode:
		act = TryLock{b}
	case LockStatMode:
		act = LockStat{b}
	}
	return act
}

type base struct {
	path     string
	duration string
}

type LockAct struct {
	base
}

func (l LockAct) Act() error {
	d, err := strconv.Atoi(l.duration)
	if err != nil {
		return err
	}
	duration := time.Minute * time.Duration(d)
	mutex := flock.New(l.path)
	if mutex.Locked() {
		return fmt.Errorf("mutex is locked! ")
	}
	if err := mutex.Lock(); err != nil {
		return err
	}
	log.Debug("start time.Sleep....")
	time.Sleep(duration)
	log.Debug("stop time.Sleep....")
	return nil
}

type TryLock struct {
	base
}

func (l TryLock) Act() error {
	d, err := strconv.Atoi(l.duration)
	if err != nil {
		return err
	}
	duration := time.Second * time.Duration(d)

	mutex := flock.New(l.path)
	for {
		time.Sleep(duration)
		ok, err := mutex.TryLock()
		if err != nil {
			return err
		}
		fmt.Println("try lock status => ", ok)
		if ok {
			break
		}
	}

	fmt.Println("try lock status ok! waiting 15s unlock! ")
	time.Sleep(time.Second * 15)
	if err = mutex.Unlock(); err != nil {
		return err
	}

	fmt.Println("try lock status unlock! waiting 5s exit! ")
	time.Sleep(time.Second * 5)
	return nil
}

type LockStat struct {
	base
}

func (l LockStat) Act() error {
	d, err := strconv.Atoi(l.duration)
	if err != nil {
		return err
	}
	duration := time.Second * time.Duration(d)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	tick := time.NewTicker(time.Second)
	mutex := flock.New(l.path)
	for {
		select {
		case <-tick.C:
			ok, err := mutex.TryRLock()
			if err != nil {
				return fmt.Errorf("mutex.status tryRLock error: %v", err)
			}
			fmt.Println("mutex.status => ", ok)
		case <-ctx.Done():
			fmt.Println("mutex.status timeout. ")
			return nil
		}
	}
}
