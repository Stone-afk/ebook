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
	// sc => SubConn
	// sci => SubConnInfo
	for sc, sci := range info.ReadySCs {
		cc := &conn{cc: sc}
		md, ok := sci.Address.Metadata.(map[string]any)
		if ok {
			weightVal := md["weight"]
			weight, _ := weightVal.(float64)
			cc.weight = int(weight)
			//group, _ := md["group"]
			//cc.group =group
			cc.labels = md["labels"].([]string)
		}

		if cc.weight == 0 {
			// 可以给个默认值
			cc.weight = 10
		}
		cc.currentWeight = cc.weight
		conns = append(conns, cc)
	}
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
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var total int
	var maxCC *conn
	// 要计算当前权重
	//label := info.Ctx.Value("label")
	for _, cc := range p.conns {
		if !cc.available {
			continue
		}
		// 如果要是 cc 里面的所有标签都不包含这个 label ，就跳过
		// 性能最好就是在 cc 上用原子操作
		// 但是筛选结果不会严格符合 WRR 算法
		// 整体效果可以
		// cc.lock.Lock()
		total += cc.weight
		cc.currentWeight = cc.currentWeight + cc.weight
		if maxCC == nil || cc.currentWeight > maxCC.currentWeight {
			maxCC = cc
		}
		//cc.lock.Unlock()
	}
	// 更新
	maxCC.currentWeight = maxCC.currentWeight - total
	// maxCC 就是挑出来的
	return balancer.PickResult{
		SubConn: maxCC.cc,
	}, nil

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
