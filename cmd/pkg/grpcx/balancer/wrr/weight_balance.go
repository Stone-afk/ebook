package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

const name = "custom_wrr"

// balancer.Balancer 接口
// balancer.Builder 接口
// balancer.Picker 接口
// base.PickerBuilder 接口
// 你可以认为，Balancer 是 Picker 的装饰器
func init() {

}

// 传统版本的基于权重的负载均衡算法

type PickerBuilder struct{}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*conn, 0, len(info.ReadySCs))
	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	//	 这个才是真的执行负载均衡的地方
	conns []*conn
	mutex sync.Mutex
}

// Pick 在这里实现基于权重的负载均衡算法
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	panic("")
}

func (p *Picker) healthCheck(cc *conn) bool {
	// 调用 grpc 内置的那个 health check 接口
	return true
}

// conn 代表节点
type conn struct {
	// （初始）权重
	weight int
	labels []string
	// 有效权重
	//efficientWeight int
	currentWeight int
	available     bool
	// 假如有 vip 或者非 vip
	group string
	//	真正的，grpc 里面的代表一个节点的表达
	cc balancer.SubConn
}
