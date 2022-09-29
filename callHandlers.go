package main

import (
	"sync"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type Mod struct {
	sync.Mutex
	mode     string
	input    string
	lang     string
	fromLang string
}

func (mod *Mod) SayMod() {
	mod.Lock()
	defer mod.Unlock()
	mod.mode = "Say"
}

func (mod *Mod) Send(input string, lang string, fromLang string) {
	mod.Lock()
	defer mod.Unlock()
	mod.lang = lang
	mod.fromLang = fromLang
	mod.input = Serv(mod.mode, input, mod.lang, mod.fromLang)
}

func (mod *Mod) DiscussMod() {
	mod.Lock()
	defer mod.Unlock()
	mod.mode = "Discuss"
}

func (mod *Mod) GetInput() string {
	mod.Lock()
	defer mod.Unlock()
	return mod.input
}
