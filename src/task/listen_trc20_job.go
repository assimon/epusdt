package task

import (
	"github.com/assimon/luuu/model/data"
	"github.com/assimon/luuu/model/service"
	"github.com/assimon/luuu/util/log"
	"sync"
)

type ListenTrc20Job struct {
}

var gListenTrc20JobLock sync.Mutex

func (r ListenTrc20Job) Run() {
	gListenTrc20JobLock.Lock()
	defer gListenTrc20JobLock.Unlock()
	walletAddress, err := data.GetAvailableWalletAddress()
	if err != nil {
		log.Sugar.Error(err)
		return
	}
	if len(walletAddress) <= 0 {
		return
	}
	var wg sync.WaitGroup
	for _, address := range walletAddress {
		wg.Add(1)
		go service.Trc20CallBack(address.Token, &wg)
	}
	wg.Wait()
}
