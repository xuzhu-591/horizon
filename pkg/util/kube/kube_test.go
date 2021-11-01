package kube

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	fakek8s "k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"

	"g.hz.netease.com/horizon/pkg/util/kube/fake"
	"g.hz.netease.com/horizon/pkg/util/log"
)

func TestGetEvents(t *testing.T) {
	ctx := log.WithContext(context.Background(), "GetEvents")
	clientset := fakek8s.NewSimpleClientset()
	_, _ = clientset.CoreV1().Events("ns").Create(ctx, &v1core.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "alice",
			Namespace: "ns",
		},
		// InvolvedObject 无论是否匹配，faked clientset 都会返回
		InvolvedObject: v1core.ObjectReference{
			Kind: "Pod",
			UID:  "111111",
		},
	}, metav1.CreateOptions{})

	_, _ = clientset.CoreV1().Events("ns").Create(ctx, &v1core.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bob",
			Namespace: "ns",
		},
		InvolvedObject: v1core.ObjectReference{
			Kind: "Pod",
			UID:  "22222",
		},
	}, metav1.CreateOptions{})

	events, err := GetEvents(ctx, clientset, "ns")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(events))
}

func TestGetPods(t *testing.T) {
	ctx := log.WithContext(context.Background(), "TestGetPods")
	clientset := fakek8s.NewSimpleClientset()
	_, _ = clientset.CoreV1().Pods("ns").Create(ctx, &v1core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "alice",
			Namespace: "ns",
			Labels: map[string]string{
				"age": "10",
			},
		},
	}, metav1.CreateOptions{})

	_, _ = clientset.CoreV1().Pods("ns").Create(ctx, &v1core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bob",
			Namespace: "ns",
			Labels: map[string]string{
				"age": "10",
			},
		},
	}, metav1.CreateOptions{})

	labelSelector := fields.ParseSelectorOrDie("age=10")
	pods, err := GetPods(ctx, clientset, "ns", labelSelector.String())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pods))
}

func TestBuildClient(t *testing.T) {
	// 创建临时文件
	data := `
apiVersion: v1
clusters:
- cluster:
    server: https://kubernetes.docker.internal:6443
  name: docker-desktop
contexts:
- context:
    cluster: docker-desktop
    user: docker-desktop
  name: docker-desktop
current-context: docker-desktop
kind: Config
preferences: {}
`
	tempDir, err := ioutil.TempDir("/tmp", "fake-clientset")
	assert.Nil(t, err)
	defer cleanup(t, tempDir)

	filePath := filepath.Join(tempDir, "kube.config")
	err = ioutil.WriteFile(filePath, []byte(data), 0644)
	assert.Nil(t, err)

	// 构建 client set
	_, _, err = BuildClient(filePath)
	assert.Nil(t, err)
}

func TestExec(t *testing.T) {
	// 创建 restClient
	restClient := fake.NewFakeClient()

	// 创建空的 config
	config := fake.NewEmptyClient()

	// 设置 clientset 底层使用的 http client
	clientset := kubernetes.NewForConfigOrDie(config)
	clientset.CoreV1().RESTClient().(*restclient.RESTClient).Client = restClient.Client

	// 执行 kubectl exec 动作，它包括两个请求
	// GET  'https://kubernetes.docker.internal:6443/api/v1/namespaces/default/pods/kube-proxy'
	// POST 'https://kubernetes.docker.internal:6443/api/v1/namespaces/default/pods/kube-proxy/exec
	containerRef := ContainerRef{
		Config:        config,
		KubeClientset: clientset,
		Namespace:     "test",
		Pod:           "foo1",
		Container:     "bar",
	}
	ctx := log.WithContext(context.Background(), "TestExec")
	stdout, stderr, err := Exec(ctx, containerRef, []string{"ls"}, &fake.RemoteExecutor{Client: restClient.Client})
	assert.Nil(t, err)
	assert.Equal(t, stdout,
		"http://localhost/api/v1/namespaces/test/pods/foo1/exec?command=ls&container=bar&stderr=true&stdout=true")
	assert.Equal(t, stderr, "")
}

func cleanup(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Failed to clean up %v: %v", path, err)
	}
}