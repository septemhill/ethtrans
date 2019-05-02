package main

type ServiceManager struct {
	servs []Service
}

func (sm *ServiceManager) StartServices() {
	for i := 0; i < len(sm.servs); i++ {
		sm.servs[i].Start()
	}
}

func (sm *ServiceManager) StopServices() {
	for i := 0; i < len(sm.servs); i++ {
		sm.servs[i].Stop()
	}
}

func (sm *ServiceManager) AddServices(servs ...Service) {
	for _, serv := range servs {
		sm.servs = append(sm.servs, serv)
	}
}

func NewServiceManager() *ServiceManager {
	sm := &ServiceManager{
		servs: make([]Service, 0),
	}

	return sm
}
