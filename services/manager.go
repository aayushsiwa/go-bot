package services

type Manager struct {
	workers []Worker
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Add(w Worker) {
	m.workers = append(m.workers, w)
}

func (m *Manager) StartAll() {
	for _, w := range m.workers {
		w.Start()
	}
}

func (m *Manager) StopAll() {
	for _, w := range m.workers {
		w.Stop()
	}
}
