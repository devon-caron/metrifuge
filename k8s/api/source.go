package api

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/devon-caron/metrifuge/global"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type Source interface {
	GetSourceInfo() string

	// StartLogStream starts a log stream for the source
	// kClient is the kubernetes client
	// nonK8sConfig is the non-kubernetes config
	// stopCh is the channel to signal the end of the log stream
	// This function assumes k8s is active until the rest config is checked. If the k8s config is not present, it will use the non-k8s config.
	StartLogStream(kClient *K8sClientWrapper, nonK8sConfig map[string]interface{}, stopCh <-chan struct{}) error
	GetNewLogs() []string
}

type PVCSource struct {
	PVC struct {
		Name string `json:"name" yaml:"name"`
	} `json:"pvc" yaml:"pvc"`
	LogFilePath string `json:"logFilePath" yaml:"logFilePath"`
}

type PodSource struct {
	Pod    Pod `json:"pod" yaml:"pod"`
	stream io.ReadCloser
	buffer []string
}

type Pod struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Container string `json:"container" yaml:"container"`
}

func (pvc *PVCSource) GetSourceInfo() string {
	return fmt.Sprintf("PVC: %s, Log File Path: %s", pvc.PVC.Name, pvc.LogFilePath)
}

func (pvc *PVCSource) StartLogStream(kClient *K8sClientWrapper, nonK8sConfig map[string]interface{}, stopCh <-chan struct{}) error {
	// may need to implement mount sockets for this to work
	return nil
}

func (pvc *PVCSource) GetNewLogs() []string {
	return nil
}

func (pod *PodSource) GetSourceInfo() string {
	return fmt.Sprintf("Pod: %s, Container: %s, Namespace: %s", pod.Pod.Name, pod.Pod.Container, pod.Pod.Namespace)
}

func (pod *PodSource) StartLogStream(kClient *K8sClientWrapper, nonK8sConfig map[string]interface{}, stopCh <-chan struct{}) error {
	if kClient == nil {
		panic("kClient is nil, nonK8sConfig must be provided")
	}
	stopChContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-stopCh
		cancel()
	}()

	logrus.Infof("attempting to start log stream for pod: %v", pod.GetSourceInfo())

	var stream io.ReadCloser
	maxRetries, err := strconv.Atoi(global.LOG_SOURCE_RETRIES)
	if err != nil {
		return fmt.Errorf("failed to convert log source retries to int: %v", err)
	}
	delay, err := strconv.Atoi(global.LOG_SOURCE_DELAY)
	if err != nil {
		return fmt.Errorf("failed to convert log source delay to int: %v", err)
	}
	for retries := 0; retries < maxRetries; retries++ {
		logrus.Errorf("attempting to start log stream for pod, attempt %v: %v", retries+1, pod.GetSourceInfo())

		stream, err = kClient.Clientset().CoreV1().Pods(pod.Pod.Namespace).GetLogs(pod.Pod.Name, &v1.PodLogOptions{
			Container: pod.Pod.Container,
			Follow:    true,
		}).Stream(stopChContext)
		if err != nil {
			logrus.Errorf("failed to get log stream: %v", err)
			time.Sleep(time.Duration(delay) * time.Second)
			continue
		}
		pod.stream = stream
		break
	}
	if err != nil {
		return fmt.Errorf("failed to get log stream: %v", err)
	}

	if stream == nil {
		return fmt.Errorf("log stream is nil")
	}

	// Create a scanner to read line by line
	scanner := bufio.NewScanner(stream)

	debugCounter := 0
	for scanner.Scan() {
		logLine := scanner.Text()
		pod.buffer = append(pod.buffer, logLine)
		debugCounter++
		if debugCounter > 10 {
			logrus.Infof("received 10 logs from pod: %v", pod.GetSourceInfo())
			debugCounter = 0
		}
	}

	logrus.Infof("finished reading logs from pod: %v", pod.GetSourceInfo())

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (pod *PodSource) GetNewLogs() []string {
	newLogs := make([]string, len(pod.buffer))
	copy(newLogs, pod.buffer)
	pod.buffer = make([]string, 0)
	return newLogs
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

func (locs *LocalSource) StartLogStream(kClient *K8sClientWrapper, nonK8sConfig map[string]interface{}, stopCh <-chan struct{}) error {
	return nil
}

func (locs *LocalSource) GetNewLogs() []string {
	return nil
}

func (cs *CmdSource) GetSourceInfo() string {
	return fmt.Sprintf("Command: %s", cs.Command)
}

func (cs *CmdSource) StartLogStream(kClient *K8sClientWrapper, nonK8sConfig map[string]interface{}, stopCh <-chan struct{}) error {
	return nil
}

func (cs *CmdSource) GetNewLogs() []string {
	return nil
}
