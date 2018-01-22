package bll

// 球信息 到 球网络信息

import (
	"usercmd"
)

var (
	// XXX 清理重复定义
	MsgPosScaleRate float64 = 100 // 服务器内部坐标 * MsgPosScaleRate = 发送给客户端的坐标
)

func FoodToMsgBall(ball *BallFood) *usercmd.MsgBall {
	return &usercmd.MsgBall{
		Id:   ball.id,
		Type: int32(ball.typeID),
		X:    int32(ball.Pos.X * MsgPosScaleRate),
		Y:    int32(ball.Pos.Y * MsgPosScaleRate),
	}
}

func FeedToMsgBall(ball *BallFeed) *usercmd.MsgBall {
	cmd := &usercmd.MsgBall{
		Id:   ball.id,
		Type: int32(ball.typeID),
		X:    int32(ball.Pos.X * MsgPosScaleRate),
		Y:    int32(ball.Pos.Y * MsgPosScaleRate),
	}
	return cmd
}

func SkillToMsgBall(ball *BallSkill) *usercmd.MsgBall {
	cmd := &usercmd.MsgBall{
		Id:   ball.id,
		Type: int32(ball.BallType),
		X:    int32(ball.Pos.X * MsgPosScaleRate),
		Y:    int32(ball.Pos.Y * MsgPosScaleRate),
	}
	return cmd
}

func PlayerBallToMsgBall(ball *BallPlayer) *usercmd.MsgPlayerBall {
	cmd := &usercmd.MsgPlayerBall{
		Id:    ball.id,
		Level: uint32(ball.GetAnimalId()),
		Hp:    uint32(ball.GetHP()),
		Mp:    uint32(ball.GetMP()),
		X:     int32(ball.Pos.X * MsgPosScaleRate),
		Y:     int32(ball.Pos.Y * MsgPosScaleRate),
		Angle: int32(ball.player.GetAngle()),
		Face:  uint32(ball.player.GetFace()),
	}

	return cmd
}
