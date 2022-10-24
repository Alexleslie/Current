package geecache

import pb "Current/Cdn/geecache/geecachepb"

/*
原流程
						  	  是
	接收key --> 检查是否被缓存 -----> 返回缓存值（1）
					|   否                        是
					| -----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值（2）
							    	|   否
									| -----> 调用”回调函数“，获取值并添加到缓存 --> 返回缓存值（3）

使用一致性哈希算法选择节点
						  是								  是
	|-----> 是否是远程节点 -----> HTTP客户端访问远程节点 --> 成功？ -----> 服务端返回返回值
				| 否                                      | 否
				|-------------------------------------------->回退到本地节点处理
*/

// PeerPicker 客户端选择，根据一致哈希算法选择节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool) //根据传入的key选择相应节点PeerGetter
}

// PeerGetter 客户端获取缓存信息
type PeerGetter interface {
	GetValueFromRemotePeer(in *pb.Request, out *pb.Response) error //从对应group查找缓存值（PeerGetter对应与上述流程的HTTP客户端）
}
