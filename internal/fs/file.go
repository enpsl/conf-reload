// Copyright 2023 enpsl. All rights reserved.

// Package fs is Instantiation of broker interface
// File change notification based on fsnotify
// conf-reload and its internal packages.

package fs

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/enpsl/conf-reload/internal/base"
	"github.com/enpsl/conf-reload/internal/errors"
	"github.com/enpsl/conf-reload/internal/log"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

type unmarshaller func([]byte, interface{}) error

type FsBroker struct {
	logger       *log.Logger
	notifyCh     chan struct{}
	once         sync.Once
	unmarshaller unmarshaller
	dir          string
	abs          string
	ext          FileExtType
	wg           sync.WaitGroup
}

type FileExtType string

const FileExtToml FileExtType = "toml"
const FileExtJson FileExtType = "json"
const FileExtYaml FileExtType = "yaml"

var UnmarshallerMap = map[FileExtType]unmarshaller{
	FileExtToml: toml.Unmarshal,
	FileExtJson: json.Unmarshal,
	FileExtYaml: yaml.Unmarshal,
}

func ExtParser(file string) FileExtType {
	ext := filepath.Ext(file)
	switch ext {
	case ".toml":
		return FileExtToml
	case ".json":
		return FileExtJson
	case ".yaml":
		return FileExtYaml
	case ".yml":
		return FileExtYaml
	}
	return ""
}

func NewFs(path string, logger *log.Logger) (error, *FsBroker) {
	fs := new(FsBroker)
	fs.notifyCh = make(chan struct{})

	abs, err := filepath.Abs(path)

	if err != nil {
		return errors.ErrFormat(errors.ErrInvalidFilePath, err), nil
	}

	fs.abs = abs

	dir, err := base.FindParentDir(abs)
	if err != nil {
		return errors.ErrFormat(errors.ErrInvalidFilePath, err), nil
	}

	ext_type := ExtParser(abs)
	if _, ok := UnmarshallerMap[ext_type]; !ok {
		return errors.ErrFormat(errors.ErrInvalidFileExt, errors.New("ext is unsupport")), nil
	}
	fs.unmarshaller = UnmarshallerMap[ext_type]
	fs.dir = dir
	fs.logger = logger
	fs.ext = ext_type
	return nil, fs
}

func (fs *FsBroker) Parse(content []byte) (error, map[string]interface{}) {
	var config = make(map[string]interface{})
	err := fs.unmarshaller(content, &config)
	if err != nil {
		return errors.ErrFormat(errors.ErrUnmarshaller, err), nil
	}
	return nil, config
}

func (fs *FsBroker) LoadContent() ([]byte, error) {
	return os.ReadFile(fs.abs)
}

func (fs *FsBroker) Watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		fs.logger.Fatalf("new file watcher error:%w", err.Error())
	}
	defer w.Close()

	configFile := filepath.Clean(fs.abs)
	realConfigFile, _ := filepath.EvalSymlinks(fs.abs)

	fs.wg.Add(1)
	go func() {
		defer fs.wg.Done()
		for {
			select {
			case event := <-w.Events:
				// Compatible with soft links
				currentConfigFile, _ := filepath.EvalSymlinks(fs.abs)
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if (filepath.Clean(event.Name) == configFile && event.Op&writeOrCreateMask != 0) ||
					(currentConfigFile != "" && currentConfigFile != realConfigFile) {
					realConfigFile = currentConfigFile
					fs.logger.Debugf("modified file:%s, %s", event.Name, realConfigFile)
					fs.notifyCh <- struct{}{}
				}
			case err := <-w.Errors:
				fs.logger.Errorf("read watch error:" + err.Error())
			}
		}
	}()
	err = w.Add(fs.dir)
	if err != nil {
		fs.logger.Fatal(err)
	}
	fs.wg.Wait()
}

func (fs *FsBroker) Decode(input interface{}, output interface{}, weaklyTypedInput bool) error {
	config := mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		Result:           output,
		TagName:          string(fs.ext),
		WeaklyTypedInput: weaklyTypedInput,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return errors.ErrFormat(errors.ErrBrokerDecode, fmt.Errorf("decode err %w", err))
	}
	return decoder.Decode(input)
}

func (fs *FsBroker) Notify() <-chan struct{} {
	return fs.notifyCh
}

func (fs *FsBroker) Close() error {
	fs.once.Do(func() {
		close(fs.notifyCh)
	})
	return nil
}
