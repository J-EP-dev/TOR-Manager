package executor

import "tormanager/pkg/instance"

type TorManager struct {
	instances map[*instance.TorInstance]bool
	addInstance chan *instance.TorInstance
	removeInstance chan *instance.TorInstance
}

func (manager *TorManager) Start()  {
	for {
		select {
		case i := <- manager.addInstance:
			manager.instances[i] = true
		case i := <- manager.removeInstance:
			delete(manager.instances, i)
		}
	}
}

func (manager *TorManager) InitInstance(instance *instance.TorInstance) {
	manager.addInstance <- instance
	instance.Run()
	manager.removeInstance <- instance
}

func (manager *TorManager) GetInstances() map[*instance.TorInstance]bool {
	return manager.instances
}
