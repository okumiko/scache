package scache

import pb "scache/scachepb"

type PeerPicker interface {
	//如果是远端key，返回true和获取远端的ProtorGetter接口
	//如果是本地key，返回false，直接返回本地缓存
	PickPeer(key string) (peer ProtoGetter, ok bool)
}

//通过pb协议传输节点间的数据
type ProtoGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
