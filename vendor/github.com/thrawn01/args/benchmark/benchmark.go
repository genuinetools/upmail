package benchmark

import (
	"fmt"
	"net/http"
	"sync"
)

type OptValues map[string]string
type Options struct {
	values OptValues
}

func (self *Options) Get(key string) string {
	return self.values[key]
}

type ArgParser interface {
	SetOpts(values OptValues)
	GetOpts() *Options
}

type Api struct {
	Parser ArgParser
}

func (self *Api) NewServer() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", self.Index)
	return router
}

func (self *Api) Index(resp http.ResponseWriter, req *http.Request) {
	opts := self.Parser.GetOpts()
	fmt.Fprintf(resp, "TestValue: %s", opts.Get("test-value"))
}

// No Thread Safety
type ParserBench struct {
	opts *Options
}

// =================================================================
func NewParserBench(values OptValues) *ParserBench {
	parser := &ParserBench{&Options{}}
	parser.SetOpts(values)
	return parser
}

func (self *ParserBench) SetOpts(values OptValues) {
	self.opts = &Options{values}
}

func (self *ParserBench) GetOpts() *Options {
	return self.opts
}

type ParserBenchMutex struct {
	opts  *Options
	mutex *sync.Mutex
}

// =================================================================
func NewParserBenchMutex(values OptValues) *ParserBenchMutex {
	parser := &ParserBenchMutex{&Options{}, &sync.Mutex{}}
	parser.SetOpts(values)
	return parser
}

func (self *ParserBenchMutex) SetOpts(values OptValues) {
	self.mutex.Lock()
	self.opts = &Options{values}
	self.mutex.Unlock()
}

func (self *ParserBenchMutex) GetOpts() *Options {
	self.mutex.Lock()
	defer func() {
		self.mutex.Unlock()
	}()
	return self.opts
}

// =================================================================
type ParserBenchRWMutex struct {
	opts  *Options
	mutex *sync.RWMutex
}

func NewParserBenchRWMutex(values OptValues) *ParserBenchRWMutex {
	parser := &ParserBenchRWMutex{&Options{}, &sync.RWMutex{}}
	parser.SetOpts(values)
	return parser
}

func (self *ParserBenchRWMutex) SetOpts(values OptValues) {
	self.mutex.Lock()
	self.opts = &Options{values}
	self.mutex.Unlock()
}

func (self *ParserBenchRWMutex) GetOpts() *Options {
	self.mutex.Lock()
	defer func() {
		self.mutex.Unlock()
	}()
	return self.opts
}

// =================================================================
type ParserBenchChannel struct {
	opts *Options
	get  chan *Options
	set  chan *Options
	done chan bool
}

func NewParserBenchChannel(values OptValues) *ParserBenchChannel {
	parser := &ParserBenchChannel{&Options{}, make(chan *Options), make(chan *Options), make(chan bool)}
	parser.Open()
	parser.SetOpts(values)
	return parser
}

func (self *ParserBenchChannel) Open() {
	go func() {
		defer func() {
			close(self.get)
			close(self.set)
			close(self.done)
		}()
		for {
			select {
			case self.get <- self.opts:
			case value := <-self.set:
				self.opts = value
			case <-self.done:
				return
			}
		}
	}()
}

func (self *ParserBenchChannel) Close() {
	self.done <- true
}

func (self *ParserBenchChannel) SetOpts(values OptValues) {
	self.set <- &Options{values}
}

func (self *ParserBenchChannel) GetOpts() *Options {
	return <-self.get
}
