/**
 * @Author: zzy
 * @Email: zhangzhongyuan@didiglobal.com
 * @Description:
 * @File: dns.go
 * @Package: dns
 * @Version: 1.0.0
 * @Date: 2023/12/27 02:39
 */

package xdns

import (
	_ "net"
	_ "unsafe"

	"os"
	"sync"
	"sync/atomic"
	"time"
)

// A resolverConfig represents a DNS stub resolver configuration.
type resolverConfig struct {
	initOnce sync.Once // guards init of resolverConfig

	// ch is used as a semaphore that only allows one lookup at a
	// time to recheck resolv.conf.
	ch          chan struct{} // guards lastChecked and modTime
	lastChecked time.Time     // last time resolv.conf was checked

	dnsConfig atomic.Pointer[struct{}] // parsed resolv.conf structure used in lookups
}

//go:linkname resolvConf net.resolvConf
var resolvConf resolverConfig

//go:linkname dnsReadConfig net.dnsReadConfig
func dnsReadConfig(filename string) *struct{}

func Hook(filename string) (err error) {
	if _, err = os.Stat(filename); err != nil {
		return
	}

	resolvConf.dnsConfig.Store(dnsReadConfig(filename))

	return
}

// TODO 后续考虑移动到 golib 库中，也把miskit的dns移动到 golib 库中
// TODO ddns 项目中有完整的实现，miskit中的Hook不算是Hook，就是dns服务器
