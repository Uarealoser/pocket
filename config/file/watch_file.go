package file

import (
	"Uarealoser/pocket/config"
	"Uarealoser/pocket/config/ini"
	"os"
	"sync"
	"time"
)

var updateInterval = time.Second * 10

type FileWatcher struct {
	config.Watcher
	LastModifyTime time.Time
	Path           string
	Priority       int
	WatchErrorFunc config.WatchErrorFunc
	Stopped        bool
}

func (f *FileWatcher) SetPriority(priority int) {
	f.Priority = priority
}

func (f *FileWatcher) GetPriority() int {
	return f.Priority
}

func (f *FileWatcher) Init(notifyChan chan *config.ConfigNotify) error {
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
	notifyChan <- &config.ConfigNotify{Watcher: f, Config: cfg, Wg: wg}
	wg.Wait()
	return nil
}

// Watch _
func (f *FileWatcher) Watch(notifyChan chan *config.ConfigNotify) {

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
				config.DefaultWatchError(err)
			}
			continue
		}

		// 修改时间不匹配,就执行更新配置
		if !f.LastModifyTime.Equal(stat.ModTime()) {
			cfg, err := f.read()
			if err == nil {
				notifyChan <- &config.ConfigNotify{Watcher: f, Config: cfg}
				f.LastModifyTime = stat.ModTime()
			} else {
				if f.WatchErrorFunc != nil {
					f.WatchErrorFunc(err)
				} else {
					config.DefaultWatchError(err)
				}
			}
		}
	}
}

func (f *FileWatcher) read() (config.Config, error) {
	return ini.NewIniConfigFromFile(f.Path)
}

func (f *FileWatcher) SetWatchErrorFunc(errorFunc config.WatchErrorFunc) {
	f.WatchErrorFunc = errorFunc
}

func (f *FileWatcher) Stop() {
	f.Stopped = true
}
