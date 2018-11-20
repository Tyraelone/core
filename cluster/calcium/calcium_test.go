package calcium

import (
	networkmocks "github.com/projecteru2/core/network/mocks"
	schedulermocks "github.com/projecteru2/core/scheduler/mocks"
	sourcemocks "github.com/projecteru2/core/source/mocks"
	storemocks "github.com/projecteru2/core/store/mocks"
	"github.com/projecteru2/core/types"
)

func NewTestCluster() *Calcium {
	c := &Calcium{}
	c.config = types.Config{}
	c.store = &storemocks.Store{}
	c.scheduler = &schedulermocks.Scheduler{}
	c.network = &networkmocks.Network{}
	c.source = &sourcemocks.Source{}
	return c
}