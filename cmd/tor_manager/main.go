package main

import (
	"fmt"
	"time"
	"tormanager/pkg/executor"
	"tormanager/pkg/instance"
)

func main() {
	// TODO Struct para instanciar la config
	// TODO Detectar códigos de error distintos a 0
	// TODO Detectar cuando falla
	// TODO Detectar puertos duplicados antes de iniciar
	// DONE Detectar error al bindear puertos.
	// TODO Struct de excepciones.
	torExecutor := executor.NewExecutor()
	torExecutor.Start()

	newInstance := instance.NewTorInstance(65535)
	newInstance2 := instance.NewTorInstance(65534)
	newInstance3 := instance.NewTorInstance(65533)

	torExecutor.StartInstance(newInstance)
	torExecutor.StartInstance(newInstance2)
	torExecutor.StartInstance(newInstance3)
	newInstance.WaitForStart()
	newInstance2.WaitForStart()
	newInstance3.WaitForStart()

	fmt.Println(newInstance.IsRunning())
	fmt.Println(newInstance2.IsRunning())
	fmt.Println(newInstance3.IsRunning())

	time.Sleep(10 * time.Second)

	torExecutor.StopInstance(newInstance)
	torExecutor.StopInstance(newInstance2)
	torExecutor.StopInstance(newInstance3)
	newInstance.WaitForStop()
	newInstance2.WaitForStop()
	newInstance3.WaitForStop()
}
