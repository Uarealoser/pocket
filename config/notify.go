package config

import "sync"

type ConfigNotify struct {
	Watcher
	Config
	Wg *sync.WaitGroup
}
