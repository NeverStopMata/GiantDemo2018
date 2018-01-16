package bll

// 技能球

import (
	bmath "base/math"
	"roomserver/game/interfaces"
	"roomserver/util"
	"usercmd"
)

type BallSkill struct {
	BallMove
	player IScenePlayer
	Skill  interfaces.ISkillBall
}

func NewBallSkill(_ballType usercmd.BallType, id uint32, x, y, radius float64, player IScenePlayer) *BallSkill {
	ball := BallSkill{
		BallMove: BallMove{
			BallFood: BallFood{
				id:       id,
				typeID:   uint16(_ballType),
				BallType: _ballType,
				Pos:      util.Vector2{float64(x), float64(y)},
				radius:   float64(radius),
			},
		},
		player: player,
	}

	// ball.Skill = MyProvider.NewSkillBall(player, &ball)
	ball.Skill = player.NewSkillBall(&ball) // XXX 临时实现，应该有更好的方法
	return &ball
}

func (this *BallSkill) Move(pertime float64, scene IScene) bool {
	if this.speed.IsEmpty() == false {
		if this.player.GetFrame() >= this.Skill.GetBeginFrame() {
			if false {
				pos := this.PhysicObj.GetPostion()
				this.Pos = util.Vector2{float64(pos.X), float64(pos.Y)}
				this.PhysicObj.SetVelocity(&bmath.Vector2{float32(this.speed.X), float32(this.speed.Y)})
			}

			this.Pos.X = this.Pos.X + this.speed.X/2
			this.Pos.Y = this.Pos.Y + this.speed.Y/2
		}
	}
	return true
}
