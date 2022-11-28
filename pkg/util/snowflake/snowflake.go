package snowflake

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	numberBits uint8 = 12 // 表示每个集群下的每个节点，1毫秒内可生成的id序号的二进制位 对应上图中的最后一段
	workerBits uint8 = 10 // 每台机器(节点)的ID位数 10位最大可以有2^10=1024个节点数 即每毫秒可生成 2^12-1=4096个唯一ID 对应上图中的倒数第二段
	// 这里求最大值使用了位运算，-1 的二进制表示为 1 的补码，感兴趣的同学可以自己算算试试 -1 ^ (-1 << nodeBits) 这里是不是等于 1023
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点ID的最大值，用于防止溢出
	numberMax   int64 = -1 ^ (-1 << numberBits) // 同上，用来表示生成id序号的最大值
	timeShift         = workerBits + numberBits // 时间戳向左的偏移量
	workerShift       = numberBits              // 节点ID向左的偏移量
	epoch       int64 = 1668354334000
)

var worker *Worker

func init() {
	if newWorker, err := NewWorker(0); err != nil {
		println(err)
		panic("new snowflake error")
	} else {
		worker = newWorker
	}
}

// Worker 定义一个woker工作节点所需要的基本参数
type Worker struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	timestamp int64      // 记录上一次生成id的时间戳
	workerId  int64      // 该节点的ID
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

// NewWorker 实例化一个工作节点  workerId 为当前节点的id
func NewWorker(workerId int64) (*Worker, error) {
	// 要先检测workerId是否在上面定义的范围内
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Worker) GetId() int64 {
	// 获取id最关键的一点 加锁 加锁 加锁
	w.mu.Lock()
	defer w.mu.Unlock() // 生成完成后记得 解锁 解锁 解锁

	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if w.timestamp == now {
		w.number++

		// 这里要判断，当前工作节点是否在1毫秒内已经生成numberMax个ID
		if w.number > numberMax {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		w.number = 0
		// 下面这段代码看到很多前辈都写在if外面，无论节点上次生成id的时间戳与当前时间是否相同 都重新赋值  这样会增加一丢丢的额外开销 所以我这里是选择放在else里面
		w.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	}
	ID := int64((now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}

func ipToInt(ip string) int64 {
	bits := strings.Split(ip, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func getIntranetIp() string {
	var (
		ip = "127.0.0.1"
	)
	address, err := net.InterfaceAddrs()

	if err != nil {
		return ip
	}
	for _, address := range address {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}
	return ip
}

func NextId() string {
	return strconv.FormatInt(worker.GetId(), 10)
}
