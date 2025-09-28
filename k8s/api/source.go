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
