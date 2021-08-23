package executor

import (
	"os/exec"
	"tormanager/pkg/instance"
)

type Executor struct {
	manager TorManager
}

func NewExecutor() Executor {
	return Executor{}
}

func (e *Executor) Start() {
	_, err := exec.LookPath("tor")
	if err != nil {
		panic("Tor not installed!")
	}

	e.manager = TorManager{
		instances:      make(map[*instance.TorInstance]bool),
		addInstance:    make(chan *instance.TorInstance),
		removeInstance: make(chan *instance.TorInstance),
	}
	go e.manager.Start()
}

func (e *Executor) StartInstance(instance *instance.TorInstance) {
	go e.manager.InitInstance(instance)
}

func (e *Executor) StopInstance(instance *instance.TorInstance) {
	instance.Stop()
}

func (e *Executor) GetManager() TorManager {
	return e.manager
}
