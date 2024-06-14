package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Renders a templ component and returns an error if any
// options can be passed to the component handler to modify the behavior
func RenderTempl(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}

type CacheOptions struct {
	MaxAge               time.Duration // Default: 1 hour
	StaleWhileRevalidate time.Duration // Default: 5 minutes
	StaleIfError         time.Duration // Default: 5 minutes
}

// SetCacheHeader sets the cache headers for the response based on the provided options.
// The function checks for negative durations and sets default values as follows:
//
// - MaxAge: 1 hour if the provided value is negative.
//
// - StaleWhileRevalidate: 5 minutes if the provided value is negative.
//
// - StaleIfError: 5 minutes if the provided value is negative.
//
// This configuration means that the response will be considered fresh for the duration of MaxAge,
// and may be revalidated in the background every StaleWhileRevalidate duration.
// In case of an error, the stale response may still be used for the duration of StaleIfError.
//
// The cache control header includes both 'max-age' for standard caches and 's-maxage' for shared caches,
// ensuring compatibility with browsers and CDNs.
func SetCacheHeader(c *fiber.Ctx, options CacheOptions) {
	// Set default values if not provided
	if options.MaxAge < 0 {
		options.MaxAge = time.Hour
	}
	if options.StaleWhileRevalidate < 0 {
		options.StaleWhileRevalidate = 5 * time.Minute
	}
	if options.StaleIfError < 0 {
		options.StaleIfError = 5 * time.Minute
	}

	// Convert the duration to seconds
	maxAge := options.MaxAge / time.Second
	staleWhileRevalidate := options.StaleWhileRevalidate / time.Second
	staleIfError := options.StaleIfError / time.Second

	// Construct the cache control header
	cacheControl := fmt.Sprintf("public, max-age=%d, s-maxage=%d, stale-while-revalidate=%d, stale-if-error=%d", maxAge, maxAge, staleWhileRevalidate, staleIfError)

	// Set the cache control header
	c.Set("Cache-Control", cacheControl)
}

// Returns the trueVal if the condition is true, otherwise it returns the falseVal.
func TernaryIf[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

var Printer = message.NewPrinter(language.English)

func GetFileModTime(file string) time.Time {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return time.Time{}
	}
	return fileInfo.ModTime()
}

func Jsonify(data interface{}) string {
	stringified, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(stringified)
}

func Truncate(description string, wordLimit int) string {
	words := strings.Fields(description)
	if len(words) > wordLimit {
		return strings.Join(words[:wordLimit], " ")
	}
	return description
}

type Promise[T any] struct {
	ResultCh chan T
	ErrorCh  chan error
}

func (p Promise[T]) Wait() (T, error) {
	err := <-p.ErrorCh
	res := <-p.ResultCh

	return res, err
}

func (p Promise[T]) Close() {
	close(p.ErrorCh)
	close(p.ResultCh)
}

func Async[T any](fn func() (T, error)) Promise[T] {
	p := Promise[T]{
		ResultCh: make(chan T, 1),
		ErrorCh:  make(chan error, 1),
	}

	go func() {
		defer p.Close()

		res, err := fn()
		p.ResultCh <- res
		p.ErrorCh <- err
	}()

	return p
}

func CacheKey(key string, args ...interface{}) string {
	var result string = fmt.Sprintf("%v", key)
	if len(args) == 0 {
		return result
	}
	for _, arg := range args {
		result += fmt.Sprintf(":%v", arg)
	}
	return result
}

type ICacheStore interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte, exp time.Duration) error
}

func Remember[T any](store ICacheStore, key string, duration time.Duration, fn func() (T, error)) (T, error) {
	cached, err := store.Get(key)
	if err == nil && cached != nil && len(cached) > 0 {
		var result T
		if err := json.Unmarshal(cached, &result); err != nil {
			return result, err
		}
		return result, nil
	}

	result, err := fn()
	if err != nil {
		return result, err
	}

	computed, err := json.Marshal(result)
	if err != nil {
		return result, err
	}
	store.Set(key, computed, duration)
	return result, nil
}

func NewCacheStore() *CacheStore {
	return &CacheStore{
		kV:       make(map[string][]byte),
		expiries: make(map[string]time.Time),
	}
}

type CacheStore struct {
	kV       map[string][]byte
	expiries map[string]time.Time
}

func (c *CacheStore) Get(key string) ([]byte, error) {
	val, ok := c.kV[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	expiry, exists := c.expiries[key]
	if exists && time.Now().After(expiry) {
		c.Delete(key)
		return nil, errors.New("key expired")
	}

	return val, nil
}

func (c *CacheStore) Set(key string, val []byte, exp time.Duration) error {
	c.kV[key] = val
	if exp > 0 {
		c.expiries[key] = time.Now().Add(exp)
	} else {
		delete(c.expiries, key)
	}
	return nil
}

func (c *CacheStore) Delete(key string) error {
	delete(c.kV, key)
	delete(c.expiries, key)
	return nil
}
