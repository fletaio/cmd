package closer

import "log"

// Closer is Closer inferface
type Closer interface {
	Close()
}

// Manager handles closers
type Manager struct {
	Names   []string
	Closers []Closer
	closed  chan struct{}
}

// NewManager returns a Manager
func NewManager() *Manager {
	cm := &Manager{
		Names:   []string{},
		Closers: []Closer{},
		closed:  make(chan struct{}),
	}
	return cm
}

// RemoveAll removes all closers
func (cm *Manager) RemoveAll() {
	cm.Names = []string{}
	cm.Closers = []Closer{}
}

// Add adds a closer with a name
func (cm *Manager) Add(Name string, c Closer) {
	cm.Names = append(cm.Names, Name)
	cm.Closers = append(cm.Closers, c)
}

// CloseAll closers all closers
func (cm *Manager) CloseAll() {
	for i, c := range cm.Closers {
		log.Println("Close", cm.Names[i])
		c.Close()
	}
	cm.closed <- struct{}{}
}

// Wait waits close all
func (cm *Manager) Wait() {
	<-cm.closed
}