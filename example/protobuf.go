package example

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"skframe/app/protobuffs/Player"
)

func TestProtoBuff() {
	var gender uint64
	var uid uint64
	var coin uint64
	name := "test"
	url := "https://www.wegame.com/account/263369/287631QhgxK.png"
	address := "0.0.0.0"
	gender = 0
	previleges := ""
	uid = 123456
	coin = 1
	// 给 proto 定义的字段赋值
	playerInfo := &Player.PlayerInfo{
		Nickname:     &name,
		HeadImageUrl: &url,
		Address:      &address,
		Gender:       &gender,
		Privileges:   &previleges,
		Uid:          &uid,
		Coin:         &coin,
	}
	info, err := proto.Marshal(playerInfo)
	if err != nil {
		fmt.Println("序列化错误", err)
	}
	fmt.Println(len(info))
	return
}
