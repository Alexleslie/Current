package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/*
缓存雪崩：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。常因为缓存服务器宕机或缓存设置了相同的过期时间引起。
*/

/*
一致性哈希算法（consistent hashing）：将key映射到2^32的空间（将该空间数字俺顺时针排列首位相连形成一个环）
计算节点/机器（通常使用节点的名称、编号、IP地址）的哈希值，放在环上。
计算key的哈希值，放置在环上，顺时针寻找到的第一个节点，就是应选取的节点/机器。
数据倾斜问题：引入虚拟节点，计算虚拟节点的hash值，放置在环上；计算key的hash值在环上顺时针寻找到应选取的虚拟节点，通过字典映射找到真实节点。
虚拟节点扩充了节点的数量，解决了节点较少的情况下容易倾斜的问题，且代价非常小，只需增加字典维护真实节点和虚拟节点的映射关系即可。
*/

// Hash 采取依赖注入的方式，允许用于替换成自定义的Hash函数，也方便测试时替换，默认为crc32.ChecksumIEEE算法。
// 一般来说，哈希函数考虑两个点：一个是碰撞率，一个是性能。比如 CRC、MD5、SHA1。
// 对于缓存来说，hash之后再根据节点数量取模，因此 hash 函数的碰撞率影响并不大，而是模的大小，也就是节点的数量比较关键，这也是引入虚拟节点的原因。
// 但是缓存对性能比较敏感。CRC 即循环冗余校验，编码简单，性能高，但安全性比较差，作为缓存的 hash算法很合适。
// 对于需要完整性校验的场合，碰撞率比较关键，而性能就比较次要了。一般使用 256位的 SHA1 算法，MD5 已经不再推荐了。

// Hash maps bytes to unit32
type Hash func(data []byte) uint32

// Map 为一致性哈希算法主数据结构
type Map struct {
	hash     Hash           // Hash映射函数
	replicas int            // 虚拟节点倍数
	keys     []int          // Sorted，哈希环，包括所有虚拟节点
	hashMap  map[int]string // 虚拟节点与真实节点的映射表，键是虚拟节点的哈希值，值是真实节点的名称
}

// New 构造函数，允许自定义虚拟节点倍数和Hash函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE //循环冗余校验产生一个32bit的校验值
	}
	return m
}

// AddPhysicalAndVirtualPeer 添加真实节点/虚拟节点的方法
func (m *Map) AddPhysicalAndVirtualPeer(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ { //对每一个真实节点，对应创建m.replicas个虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //虚拟节点的名称为strconv.Itoa(i) + key,即添加编号区分虚拟节点
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys) //哈希环上的哈希值排序
}

// GetRealPeerFromKey 根据包含虚拟节点的hash key，返回对应的真实节点
func (m *Map) GetRealPeerFromKey(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key))) //计算key的哈希值
	// Binary search for appropriate replica.顺势针寻找第一个匹配的虚拟节点
	// 二分查找，从[0,len(m.keys))中取出一个值index，index是该区间中使函数f(index)为true的最小值
	// 如果无法找到该index值，则返回len(m.keys)
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	return m.hashMap[m.keys[idx%len(m.keys)]] //hashMap映射得到真实节点

}
