package calcium

import (
	"context"
	"fmt"

	enginetypes "github.com/docker/docker/api/types"
	"github.com/projecteru2/core/cluster"
	"github.com/projecteru2/core/lock"
	"github.com/projecteru2/core/types"
	log "github.com/sirupsen/logrus"
)

// Lock is lock for calcium
func (c *Calcium) Lock(ctx context.Context, name string, timeout int) (lock.DistributedLock, error) {
	lock, err := c.store.CreateLock(name, timeout)
	if err != nil {
		return nil, err
	}
	if err = lock.Lock(ctx); err != nil {
		return nil, err
	}
	return lock, nil
}

// UnlockAll unlock all locks
func (c *Calcium) UnlockAll(ctx context.Context, locks map[string]lock.DistributedLock) {
	for _, lock := range locks {
		if err := lock.Unlock(ctx); err != nil {
			log.Errorf("[UnlockAll] Unlock failed %v", err)
		}
	}
}

// LockAndGetContainers lock and get containers
func (c *Calcium) LockAndGetContainers(ctx context.Context, IDs []string) (map[string]*types.Container, map[string]enginetypes.ContainerJSON, map[string]lock.DistributedLock, error) {
	containers := map[string]*types.Container{}
	containerJSONs := map[string]enginetypes.ContainerJSON{}
	locks := map[string]lock.DistributedLock{}
	for _, ID := range IDs {
		container, containerJSON, lock, err := c.LockAndGetContainer(ctx, ID)
		if err != nil {
			c.UnlockAll(ctx, locks)
			return nil, nil, nil, err
		}
		containers[ID] = container
		containerJSONs[ID] = containerJSON
		locks[ID] = lock
	}
	return containers, containerJSONs, locks, nil
}

// LockAndGetContainer lock and get container
func (c *Calcium) LockAndGetContainer(ctx context.Context, ID string) (*types.Container, enginetypes.ContainerJSON, lock.DistributedLock, error) {
	lock, err := c.Lock(ctx, fmt.Sprintf(cluster.ContainerLock, ID), c.config.LockTimeout)
	if err != nil {
		return nil, enginetypes.ContainerJSON{}, nil, err
	}
	log.Debugf("[LockAndGetContainer] Container %s locked", ID)
	// Get container
	container, err := c.store.GetContainer(ctx, ID)
	if err != nil {
		lock.Unlock(ctx)
		return nil, enginetypes.ContainerJSON{}, nil, err
	}
	// 确保是有这个容器的
	containerJSON, err := container.Inspect(ctx)
	if err != nil {
		lock.Unlock(ctx)
		return nil, enginetypes.ContainerJSON{}, nil, err
	}
	return container, containerJSON, lock, nil
}

// LockAndGetNode lock and get node
func (c *Calcium) LockAndGetNode(ctx context.Context, podname, nodename string) (*types.Node, lock.DistributedLock, error) {
	lock, err := c.Lock(ctx, fmt.Sprintf(cluster.NodeLock, podname, nodename), c.config.LockTimeout)
	if err != nil {
		return nil, nil, err
	}
	log.Debugf("[LockAndGetNode] Node %s locked", nodename)
	// Get node
	node, err := c.GetNode(ctx, podname, nodename)
	if err != nil {
		lock.Unlock(ctx)
		return nil, nil, err
	}
	return node, lock, nil
}