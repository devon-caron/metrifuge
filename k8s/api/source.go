package api

import "fmt"

type PVCSource struct {
	PVC struct {
		Name string `json:"name" yaml:"name"`
	} `json:"pvc" yaml:"pvc"`
	LogFilePath string `json:"logFilePath" yaml:"logFilePath"`
}

type PodSource struct {
	Pod struct {
		Name      string `json:"name" yaml:"name"`
		Container string `json:"container" yaml:"container"`
	} `json:"pod" yaml:"pod"`
}

type Source interface {
	GetSourceInfo() string
	StartLogStream(stopCh <-chan struct{})
	GetNewLogs() []string
}

func (pvc *PVCSource) GetSourceInfo() string {
	return fmt.Sprintf("PVC: %s, Log File Path: %s", pvc.PVC.Name, pvc.LogFilePath)
}

func (pvc *PVCSource) StartLogStream(stopCh <-chan struct{}) {
	// may need to implement mount sockets for this to work
}

func (pvc *PVCSource) GetNewLogs() []string {
	return nil
}

func (pod *PodSource) GetSourceInfo() string {
	return fmt.Sprintf("Pod: %s, Container: %s", pod.Pod.Name, pod.Pod.Container)
}

func (pod *PodSource) StartLogStream(stopCh <-chan struct{}) {

}

func (pod *PodSource) GetNewLogs() []string {
	return nil
}

// LocalSource contains the configuration for getting logs from a local file
type LocalSource struct {
	Path string `json:"path" yaml:"path"`
}

// CmdSource contains the configuration for getting logs from a command
// TODO: implement for given pod/container
type CmdSource struct {
	Command string `json:"command" yaml:"command"`
}

func (locs *LocalSource) GetSourceInfo() string {
	return fmt.Sprintf("Local: %s", locs.Path)
}

func (locs *LocalSource) StartLogStream(stopCh <-chan struct{}) {

}

func (locs *LocalSource) GetNewLogs() []string {
	return nil
}

func (cs *CmdSource) GetSourceInfo() string {
	return fmt.Sprintf("Command: %s", cs.Command)
}

func (cs *CmdSource) StartLogStream(stopCh <-chan struct{}) {

}

func (cs *CmdSource) GetNewLogs() []string {
	return nil
}
