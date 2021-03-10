package config

import (
	"Uarealoser/pocket/config/ini"
	"os"
	"sync"
	"time"
)

var updateInterval = time.Second * 3

type FileWatcher struct {
	Watcher
	LastModifyTime time.Time
	Path           string
	Priority       int
	WatchErrorFunc WatchErrorFunc
	Stopped        bool
}

func (f *FileWatcher) SetPriority(priority int) {
	f.Priority = priority
}

func (f *FileWatcher) GetPriority() int {
	return f.Priority
}

func (f *FileWatcher) Init(notifyChan chan *ConfigNotify) error {
	stat, err := os.Stat(f.Path)
	if err != nil {
		return err
	}
	cfg, err := f.read()
	if err != nil {
		return err
	}
	f.LastModifyTime = stat.ModTime()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	notifyChan <- &ConfigNotify{Watcher: f, Config: cfg, Wg: wg}
	wg.Wait()
	return nil
}

// Watch _
func (f *FileWatcher) Watch(notifyChan chan *ConfigNotify) {

	// 监听文件的修改时间
	for {
		if f.Stopped {
			break
		}
		time.Sleep(updateInterval)
		var stat os.FileInfo
		stat, err := os.Stat(f.Path)
		if err != nil {
			if f.WatchErrorFunc != nil {
				f.WatchErrorFunc(err)
			} else {
				DefaultWatchError(err)
			}
			continue
		}

		// 修改时间不匹配,就执行更新配置
		if !f.LastModifyTime.Equal(stat.ModTime()) {
			cfg, err := f.read()
			if err == nil {
				notifyChan <- &ConfigNotify{Watcher: f, Config: cfg}
				f.LastModifyTime = stat.ModTime()
			} else {
				if f.WatchErrorFunc != nil {
					f.WatchErrorFunc(err)
				} else {
					DefaultWatchError(err)
				}
			}
		}
	}
}

func (f *FileWatcher) read() (Config, error) {
	return ini.NewIniConfigFromFile(f.Path)
}

func (f *FileWatcher) SetWatchErrorFunc(errorFunc WatchErrorFunc) {
	f.WatchErrorFunc = errorFunc
}

func (f *FileWatcher) Stop() {
	f.Stopped = true
}
