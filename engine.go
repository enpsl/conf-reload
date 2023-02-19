package conf_reload

import (
	"github.com/enpsl/conf-reload/internal/app"
	"github.com/enpsl/conf-reload/internal/base"
	"github.com/enpsl/conf-reload/internal/errors"
	"github.com/enpsl/conf-reload/internal/fs"
	"github.com/enpsl/conf-reload/internal/log"
	"github.com/spf13/cast"
	"reflect"
	"strings"
	"sync"
	"time"
)

// conf-reload Engine,Used to coordinate and manage broker
type Engine struct {
	mu               sync.RWMutex           // deepsearch apply will lock
	RawData          []byte                 // config file original data
	LevelSplit       string                 // key get split
	WeaklyTypedInput bool                   // whether to startweak type conversion
	Logger           *log.Logger            // logger instance
	LocalStorage     *base.LRUCache         // fast cache
	Configure        map[string]interface{} // original config
	Broker           base.Broker            // broker
	Capacity         int
}

type Option func(*Engine)

type Logger interface {
	// Debug logs a message at Debug level.
	Debug(args ...interface{})

	// Info logs a message at Info level.
	Info(args ...interface{})

	// Warn logs a message at Warning level.
	Warn(args ...interface{})

	// Error logs a message at Error level.
	Error(args ...interface{})

	// Fatal logs a message at Fatal level
	// and process will exit with status set to 1.
	Fatal(args ...interface{})
}

// WithLevelSplit config file separator options
func WithLevelSplit(split string) Option {
	return func(engine *Engine) {
		engine.LevelSplit = split
	}
}

// WithWeaklyTypedInput options
// Whether to start weak type conversion
// See details https://github.com/mitchellh/mapstructure/blob/main/mapstructure_examples_test.go
func WithWeaklyTypedInput(weaklyTypedInput bool) Option {
	return func(engine *Engine) {
		engine.WeaklyTypedInput = weaklyTypedInput
	}
}

// WithLogger Logger options, The logger must be implement Logger
func WithLogger(logger Logger) Option {
	return func(engine *Engine) {
		engine.Logger = log.NewLogger(logger)
	}
}

// WithLogLevel Log level options
func WithLogLevel(level int32) Option {
	return func(engine *Engine) {
		engine.Logger.SetLevel(log.Level(level))
	}
}

func WithCapacity(capacity int) Option {
	return func(engine *Engine) {
		engine.Capacity = capacity
	}
}

// NewEngine
// Engine init
func NewEngine() *Engine {
	return &Engine{
		Logger:     log.NewLogger(nil),
		LevelSplit: app.DefaultLevelSplit,
		Configure:  make(map[string]interface{}),
		Capacity:   app.DefaultCapacity,
	}
}

// Load the configuration file information and initialize the broker.
// The broker will start an additional process to receive the file change chan notification
func (e *Engine) Load(path string, opts ...Option) error {
	for _, opt := range opts {
		opt(e)
	}

	e.LocalStorage = base.CacheConstructor(e.Capacity)

	err, broker := fs.NewFs(path, e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Broker = broker

	content, err := e.Broker.LoadContent()

	if err != nil {
		e.Logger.Fatal(err)
	}

	err = e.apply(content)

	if err != nil {
		e.Logger.Fatal(err)
	}

	go func() {
		for range e.Broker.Notify() {
			if content, err := e.Broker.LoadContent(); err == nil {
				err = e.apply(content)
				if err != nil {
					e.Logger.Error(err)
				}
			}
		}
	}()
	return nil
}

// apply
// Each time the configuration file changes,
// This method will be called to delete LocalStorage and update Configure
func (e *Engine) apply(content []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.RawData = content
	err, m := e.Broker.Parse(content)
	if err != nil {
		return err
	}
	e.Configure = m
	e.LocalStorage.Flush()
	e.Logger.Debug(e.Configure)
	return nil
}

// Get the value corresponding to the key from LocalStorage.
// If it is not available, it will be found in the broker
func (e *Engine) Get(key string) interface{} {
	local, ok := e.LocalStorage.Get(key)
	if ok {
		return local
	}

	paths := strings.Split(key, e.LevelSplit)
	e.mu.RLock()
	defer e.mu.RUnlock()
	m := e.deepSearch(e.Configure, paths[:len(paths)-1]...)
	e.Logger.Debug(m)
	deep := m[paths[len(paths)-1]]
	e.LocalStorage.Put(key, deep)
	return deep
}

// deepSearch
// Copy m to a new map
// This map will continuously save the value of the latest path level map during the iterative search process
// And convert the value to map[string]interface
// If the value of map corresponding to path is not found, map[path]inerface{} will be returned

func (e *Engine) deepSearch(m map[string]interface{}, paths ...string) map[string]interface{} {
	//深度拷贝
	copym := make(map[string]interface{})
	for k, v := range m {
		copym[k] = v
	}
	defaultMap := map[string]interface{}{}
	for len(paths) > 0 {
		if i, exists := copym[paths[0]]; exists {
			v := reflect.ValueOf(i)
			if v.Kind() != reflect.Map {
				copym[paths[0]] = defaultMap
				return copym
			}
			revertM, err := cast.ToStringMapE(i)
			if err != nil {
				copym[paths[0]] = defaultMap
				return copym
			}
			copym = revertM
		} else {
			copym = make(map[string]interface{})
		}
		paths = paths[1:]
	}
	return copym
}

// DecodeToStruct
// Depends on the work of broker decode
// If Map corresponding to key is nil, will return an error
func (e *Engine) DecodeToStruct(key string, i interface{}) error {
	if key == "" {
		e.mu.RLock()
		defer e.mu.RUnlock()
		return e.Broker.Decode(e.Configure, i, e.WeaklyTypedInput)
	}
	value := e.Get(key)
	if value == nil {
		return errors.ErrFormat(errors.ErrInvalidKey, nil)
	}
	return e.Broker.Decode(value, i, e.WeaklyTypedInput)
}

// Engine.GetInt returns the value associated with the key as string type.
func (e *Engine) GetString(key string) string {
	return cast.ToString(e.Get(key))
}

// Engine.GetInt returns the value associated with the key as bool type.
func (e *Engine) GetBool(key string) bool {
	return cast.ToBool(e.Get(key))
}

// Engine.GetInt returns the value associated with the key as int type.
func (e *Engine) GetInt(key string) int {
	return cast.ToInt(e.Get(key))
}

// Engine.GetInt64 returns the value associated with the key as int64 type.
func (e *Engine) GetInt64(key string) int64 {
	return cast.ToInt64(e.Get(key))
}

// Engine.GetFloat64 returns the value associated with the key as float64 type.
func (e *Engine) GetFloat64(key string) float64 {
	return cast.ToFloat64(e.Get(key))
}

// Engine.GetTime returns the value associated with the key as time.Time type.
func (e *Engine) GetTime(key string) time.Time {
	return cast.ToTime(e.Get(key))
}

// Engine.GetDuration returns the value associated with the key as time.Duration type.
func (e *Engine) GetDuration(key string) time.Duration {
	return cast.ToDuration(e.Get(key))
}

// Engine.GetStringSlice returns the value associated with the key as []string type.
func (e *Engine) GetStringSlice(key string) []string {
	return cast.ToStringSlice(e.Get(key))
}

// Engine.GetSlice returns the value associated with the key as []interface{} type.
func (e *Engine) GetSlice(key string) []interface{} {
	return cast.ToSlice(e.Get(key))
}

// Engine.GetStringMap returns the value associated with the key as map[string]interface{} type.
func (e *Engine) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(e.Get(key))
}

// Engine.GetStringMapString returns the value associated with the key as map[string]string type.
func (e *Engine) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(e.Get(key))
}

// Engine.GetStringMapStringSlice returns the value associated with the key as map[string][]string type.
func (e *Engine) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(e.Get(key))
}
