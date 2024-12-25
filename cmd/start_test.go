package cmd

import (
	"reflect"
	"testing"

	"fardjad.com/dqlite-vip/cluster"
)

type mockWaiter struct {
	called bool
}

func (w *mockWaiter) Wait() {
	w.called = true
}

func (w *mockWaiter) Validate(t *testing.T) {
	if !w.called {
		t.Errorf("Expected Wait to be called, but it wasn't")
	}
}

type mockClusterNode struct {
	startCalled bool
	closeCalled bool
}

func (n *mockClusterNode) Validate(t *testing.T) {
	if !n.startCalled {
		t.Errorf("Expected Start to be called, but it wasn't")
	}

	if !n.closeCalled {
		t.Errorf("Expected Close to be called, but it wasn't")
	}
}

func (n *mockClusterNode) Start() error {
	n.startCalled = true
	return nil
}

func (n *mockClusterNode) IsLeader() bool {
	return false
}

func (n *mockClusterNode) SetString(key string, value string) error {
	return nil
}

func (n *mockClusterNode) GetString(key string) (string, error) {
	return "", nil
}

func (n *mockClusterNode) Close() error {
	n.closeCalled = true
	return nil
}

type mockClusterNodeFactory struct {
	dataDir     string
	bindCluster string
	bindHttp    string
	join        []string

	instance *mockClusterNode
	err      error
}

func (f *mockClusterNodeFactory) Validate(t *testing.T, expectedDataDir, expectedBindCluster, expectedBindHttp string, expectedJoin []string) {
	if f.dataDir != expectedDataDir {
		t.Errorf("Expected dataDir to be '%s', but got: %s", expectedDataDir, f.dataDir)
	}

	if f.bindCluster != expectedBindCluster {
		t.Errorf("Expected bindCluster to be '%s', but got: %s", expectedBindCluster, f.bindCluster)
	}

	if f.bindHttp != expectedBindHttp {
		t.Errorf("Expected bindHttp to be '%s', but got: %s", expectedBindHttp, f.bindHttp)
	}

	if !reflect.DeepEqual(f.join, expectedJoin) {
		t.Errorf("Expected join to be %v, but got: %v", expectedJoin, f.join)
	}

	if f.instance == nil {
		t.Errorf("Expected ClusterNode to be created, but it wasn't")
	}

	f.instance.Validate(t)
}

func (f *mockClusterNodeFactory) NewClusterNode(dataDir, bindCluster, bindHttp string, join []string) (cluster.ClusterNode, error) {
	if f.err != nil {
		return nil, f.err
	}

	f.dataDir = dataDir
	f.bindCluster = bindCluster
	f.bindHttp = bindHttp
	f.join = join

	f.instance = &mockClusterNode{}

	return f.instance, nil
}

func TestStart(t *testing.T) {
	mockWaiter := &mockWaiter{}
	defer mockWaiter.Validate(t)

	mockClusterNodeFactory := &mockClusterNodeFactory{}
	defer mockClusterNodeFactory.Validate(
		t,
		"/tmp/data",
		"127.0.0.1:9000",
		"127.0.0.1:8080",
		[]string{},
	)

	startCmd := &start{
		waiter:             mockWaiter,
		clusterNodeFactory: mockClusterNodeFactory,
	}

	cmd := startCmd.command()
	cmd.SetArgs([]string{
		"--data-dir", "/tmp/data",
		"--bind-cluster", "127.0.0.1:9000",
		"--bind-http", "127.0.0.1:8080",
	})

	if err := cmd.Execute(); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestStartWithJoin(t *testing.T) {
	mockWaiter := &mockWaiter{}
	defer mockWaiter.Validate(t)

	mockClusterNodeFactory := &mockClusterNodeFactory{}
	defer mockClusterNodeFactory.Validate(
		t,
		"/tmp/data",
		"127.0.0.1:9000",
		"127.0.0.1:8080",
		[]string{"127.0.0.1:9001"},
	)

	startCmd := &start{
		waiter:             mockWaiter,
		clusterNodeFactory: mockClusterNodeFactory,
	}

	cmd := startCmd.command()
	cmd.SetArgs([]string{
		"--data-dir", "/tmp/data",
		"--bind-cluster", "127.0.0.1:9000",
		"--bind-http", "127.0.0.1:8080",
		"--join", "127.0.0.1:9001",
	})

	if err := cmd.Execute(); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
