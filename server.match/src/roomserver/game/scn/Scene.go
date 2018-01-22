// 包 scn 有场景类。
package scn

// 场景类

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"base/ape"
	"base/glog"
	"common"
	"roomserver/client/dbclient"
	"roomserver/conf"
	"roomserver/game/bll"
	"roomserver/game/cll"
	"roomserver/game/consts"
	"roomserver/game/interfaces"
	"roomserver/game/internal/physic"
	"roomserver/game/scn/playermgr"
	"roomserver/game/scn/plr"
	tm "roomserver/game/team"
	ri "roomserver/interfaces"
	"roomserver/util"
	"usercmd"
)

const (
	MAX_TEAM_NUM      = 5
	MAX_TEAM_USER_NUM = 5 // 队伍最大人数
	CUBE_X_NUM        = 16
	CUBE_Y_NUM        = 16
)

type Cube struct {
	state          int
	remainDistance int
}
type Scene struct {
	SceneBallHelper                                   // 分配球相关的辅助类
	SceneBirthPointHelper                             // 出生点辅助类
	ScenePlayerOffline                                // 离线玩家管理
	mapConfig             *conf.MapConfig             // 地图配置
	cellNumX              int                         // 格子最大X
	cellNumY              int                         // 格子最大Y
	roomSize              float64                     // 地图大小（长、宽相等）
	cells                 []*cll.Cell                 // 所有格子
	frame                 uint32                      // 当前帧数
	room                  IRoom                       // 所在房间
	Pool                  *MsgPool                    // 球、协议等对象分配池
	Players               map[uint64]*plr.ScenePlayer // 玩家对象
	scenePhysic           *physic.ScenePhysic         // 场景物理
	cubeNum               uint32                      // cube 数量
	CubeInf               []*Cube                     // cube信息的集合
	MovingCubes           map[uint32]int32            //当前场景中正在移动的cube
}

// NewSceneAIPlayer = ai.NewSceneAIPlayer
var NewSceneAIPlayer func(room IRoom) *plr.ScenePlayer

// NewCopyPlayerAI = ai.NewCopyPlayerAI
var NewCopyPlayerAI func(room IRoom, copy_player *plr.ScenePlayer) *plr.ScenePlayer

//场景初始化
func (this *Scene) Init(room IRoom) {
	this.room = room

	this.mapConfig = conf.GetMapConfigById(this.SceneID())
	this.scenePhysic = physic.NewScenePhysic()
	this.Players = make(map[uint64]*plr.ScenePlayer)
	this.SceneBallHelper.Init(this.mapConfig.Size)
	glog.Info("mata begin to load map")
	this.LoadMap()
	for i := 0; i < this.cellNumX*this.cellNumY; i++ {
		this.cells = append(this.cells, cll.NewCell(i))
	}
	for i := 0; i < this.cellNumX*this.cellNumY; i++ {
		this.cells = append(this.cells, cll.NewCell(i))
	}
	//mata:初始化cube表，高度状态都为0，待移动高度都为0
	this.cubeNum = CUBE_X_NUM * CUBE_Y_NUM
	for i := uint32(0); i < this.cubeNum; i++ {
		this.CubeInf = append(this.CubeInf, &Cube{
			state:          0,
			remainDistance: 0,
		})
	}
	this.MovingCubes = make(map[uint32]int32)
	this.reset()
}

func (this *Scene) LoadMap() {
	this.roomSize = this.mapConfig.Size
	this.cellNumX = int(math.Ceil(this.roomSize / cll.CellWidth))
	this.cellNumY = int(math.Ceil(this.roomSize / cll.CellHeight))
	this.scenePhysic.CreateBoard(float32(this.mapConfig.Size))
	glog.Info("fucking size:", float32(this.mapConfig.Size))
	for _, v := range this.mapConfig.Nodes {
		LoadMapObjectByConfig(v, this)
		randblock := this.GetSquare(v.Px, v.Py, v.Radius)
		for index, _ := range randblock {
			this.AppendFixedPos(int(randblock[index].X), int(randblock[index].Y))
		}
	}
}

// 重置整个房间
func (this *Scene) reset() {
	for _, cell := range this.cells {
		cell.Clean()
	}

	rand.NewSource(time.Now().Unix())

	this.CreateAllBirthPoint(this)

	for _, cell := range this.cells {
		cell.ResetMsg()
	}
}

//5帧更新
func (this *Scene) Render5() {
	this.SendRoomMsg()
}

//时间片渲染
func (this *Scene) Render() {
	now := time.Now()
	nowNano := now.UnixNano()
	var d float64 = consts.FrameTime
	for i, _ := range this.MovingCubes {
		if this.MovingCubes[i] == 0 {
			delete(this.MovingCubes, i)
		} else if this.MovingCubes[i] > 0 {
			this.MovingCubes[i] -= int32(d * consts.UpDownSpeed)
			if this.MovingCubes[i] < 0 {
				this.MovingCubes[i] = 0
				this.CubeInf[i].state += 1
			}
		} else if this.MovingCubes[i] < 0 {
			this.MovingCubes[i] += int32(d * consts.UpDownSpeed)
			if this.MovingCubes[i] > 0 {
				this.MovingCubes[i] = 0
				this.CubeInf[i].state -= 1
			}
		}
	} //实时更新cube还需要移动的高度
	if this.frame%2 == 0 {
		this.scenePhysic.Tick()
	}
	for _, player := range this.Players {
		player.Update(d, nowNano, this)
	}
	for _, cell := range this.cells {
		cell.Render(this, d, nowNano)
	}
	this.RefreshBirthPoint(d, this)
}

// 发送消息
func (this *Scene) SendRoomMsg() {
	for _, player := range this.Players {
		if player.Sess != nil {
			player.SendSceneMsg()
		}
	}

	for _, cell := range this.cells {
		cell.ResetMsg()
	}

	for _, player := range this.Players {
		player.ResetMsg()
	}
}

//添加球到场景
func (this *Scene) AddBall(ball interfaces.IBall) {
	x, y := ball.GetPos()
	cell, ok := this.GetCell(x, y)
	if ok {
		cell.Add(ball)
	}
}

//删除球球
func (this *Scene) RemoveBall(ball interfaces.IBall) {
	if nil == ball {
		return
	}
	for _, cell := range this.cells {
		cell.Remove(ball.GetID(), ball.GetType())
	}
}

//通过ID删除球球
func (this *Scene) RemoveBallByID(id uint32, typ usercmd.BallType) {
	for _, cell := range this.cells {
		cell.Remove(id, typ)
	}
}

//获取区域内的所有格子
func (this *Scene) GetAreaCells(s *util.Square) (cells []*cll.Cell) {
	minX := int(math.Max(math.Floor(s.Left/cll.CellWidth), 0))
	maxX := int(math.Min(math.Floor(s.Right/cll.CellWidth), float64(this.cellNumX-1)))
	minY := int(math.Max(math.Floor(s.Bottom/cll.CellHeight), 0))
	maxY := int(math.Min(math.Floor(s.Top/cll.CellHeight), float64(this.cellNumY-1)))
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			cells = append(cells, this.cells[i*this.cellNumX+j])
		}
	}
	return
}

//根据坐标获取格子
func (this *Scene) GetCell(px, py float64) (*cll.Cell, bool) {
	idxX := int(math.Max(math.Floor(px/cll.CellWidth), 0))
	idxY := int(math.Max(math.Floor(py/cll.CellHeight), 0))
	//glog.Info("[房间] cells:", px, "---", py, ", " , "-", idxX, "-", idxY)
	if idxX < this.cellNumX && idxY < this.cellNumY {
		return this.cells[idxY*this.cellNumX+idxX], true
	}
	return nil, false
}

//获取场景玩家
func (this *Scene) GetPlayer(id uint64) (*plr.ScenePlayer, bool) {
	player, ok := this.Players[id]
	return player, ok
}

//AddRobotPlayer 添加机器人
func (this *Scene) AddRobotPlayer(teamid, teamname uint32) {
	if 0 == teamid || 0 == teamname {
		return
	}
	this.AddPlayer(nil, nil, true, nil, teamid, teamname)
}

//AddPlayer 添加玩家到场景玩家
//robotIndex 机器人索引,等于0是普通玩家
func (this *Scene) AddPlayer(playertask ri.IPlayerTask, team *tm.Team, robot bool, copy_player *plr.ScenePlayer, teamid, teamname uint32) bool {
	var scenePlayer *plr.ScenePlayer
	reconn := false

	if !robot {
		splayer, ok := this.Players[playertask.ID()]
		scenePlayer = splayer
		if ok {
			tmpname := scenePlayer.UData().TeamName
			tmpid := scenePlayer.UData().TeamId
			scenePlayer.ID = playertask.ID()
			scenePlayer.Key = playertask.Key()
			scenePlayer.SetUData(playertask.UData())
			scenePlayer.Name = playertask.Name()
			scenePlayer.IsLive = true
			scenePlayer.IsRobot = false
			scenePlayer.SetDeadTime(0)
			scenePlayer.IsActClose = false
			scenePlayer.GetScene().RemoveBall(scenePlayer.SelfAnimal) //移除
			scenePlayer.GetScene().RemoveAnimalPhysic(scenePlayer.SelfAnimal.PhysicObj)
			ball := bll.NewBallPlayer(scenePlayer, scenePlayer.BallId)
			scenePlayer.SelfAnimal = ball
			this.AddBall(scenePlayer.SelfAnimal)
			reconn = true
			scenePlayer.Reconn = true
			scenePlayer.UData().TeamName = tmpname
			scenePlayer.UData().TeamId = tmpid
			glog.Info("AddPlayer playerback ", scenePlayer.Name, "ballid:", ball.GetID())
		} else {
			//混合模式默认观战
			scenePlayer = this.room.NewScenePlayer(playertask.UData(), playertask.Name(), false)
			scenePlayer.Reconn = false
			this.Players[playertask.ID()] = scenePlayer
		}

		scenePlayer.IsRobot = false
		scenePlayer.SetExp(0)
		scenePlayer.Sess = nil
	} else {
		if copy_player == nil {
			scenePlayer = NewSceneAIPlayer(this.room)
		} else {
			scenePlayer = NewCopyPlayerAI(this.room, copy_player)
			scenePlayer.SelfAnimal.ResetAnimalInfo(copy_player.SelfAnimal)
		}
		scenePlayer.IsRobot = true
		scenePlayer.Sess = nil
		scenePlayer.UData().TeamId = teamid
		scenePlayer.UData().TeamName = teamname
		this.Players[scenePlayer.ID] = scenePlayer
	}
	if scenePlayer.OldAnimalID != 0 {
		scenePlayer.SelfAnimal.SetAnimalId(scenePlayer.OldAnimalID)
		scenePlayer.SelfAnimal.SetHpMax(consts.DefaultMaxHP)
		scenePlayer.SelfAnimal.SetHP(consts.DefaultMaxHP)
		scenePlayer.OldAnimalID = 0
	}
	scenePlayer.Sess = playertask
	scenePlayer.SetExp(this.GetOfflineExp(scenePlayer.UData().Account))
	x, y, ex := this.GetOfflinePos(scenePlayer.UData().Account)
	if ex {
		scenePlayer.SelfAnimal.SetPosV(util.Vector2{x, y})
		scenePlayer.SelfAnimal.PhysicObj.SetPx(float32(x))
		scenePlayer.SelfAnimal.PhysicObj.SetPy(float32(y))
		scenePlayer.SelfAnimal.ResetRect()
	}

	// 把玩家发给其它人
	if this.room.RoomType() == common.RoomTypeTeam && false == scenePlayer.IsRobot && false == reconn {
		this.changeTeamRobotTeamId(scenePlayer)
	}
	othermsg := &this.Pool.MsgAddPlayer
	othermsg.Player.Id = scenePlayer.ID
	othermsg.Player.Name = scenePlayer.Name
	othermsg.Player.Local = scenePlayer.GetLocation()
	othermsg.Player.IsLive = scenePlayer.IsLive
	othermsg.Player.SnapInfo = scenePlayer.GetSnapInfo()
	othermsg.Player.Curexp = /*0*/ this.GetOfflineExp(scenePlayer.UData().Account)
	othermsg.Player.BallId = scenePlayer.SelfAnimal.GetID()
	othermsg.Player.Curmp = uint32(scenePlayer.SelfAnimal.GetMP())
	othermsg.Player.Curhp = uint32(scenePlayer.SelfAnimal.GetHP())
	othermsg.Player.Animalid = uint32(scenePlayer.SelfAnimal.GetAnimalId())
	othermsg.Player.TeamName = scenePlayer.UData().TeamName

	playermgr.GetMe().Remove(scenePlayer)
	if robot && scenePlayer.Key == "" {
		scenePlayer.Key = strconv.FormatInt(int64(this.room.ID()), 10) + strconv.FormatInt(int64(scenePlayer.ID), 10)
	}
	playermgr.GetMe().Add(scenePlayer)

	scenePlayer.UpdateView(this)
	scenePlayer.UpdateViewPlayers(this)
	scenePlayer.ResetMsg()
	if robot {
		this.room.BroadcastMsg(usercmd.MsgTypeCmd_AddPlayer, othermsg)
		return true
	}

	// 发送MsgTop消息给玩家(主要是更新EndTime)
	this.Pool.MsgTopRank.Reset()
	ltime := this.room.EndTime() - time.Now().Unix()
	if ltime > 0 {
		this.Pool.MsgTopRank.EndTime = uint32(ltime)
	} else {
		this.Pool.MsgTopRank.EndTime = 0
	}
	scenePlayer.Sess.SendCmd(usercmd.MsgTypeCmd_Top, &this.Pool.MsgTopRank)

	// 把当前场景的人、球都发给玩家
	var others []*usercmd.MsgPlayer
	var balls []*usercmd.MsgBall
	var playerballs []*usercmd.MsgPlayerBall
	for _, player := range this.Players {
		others = append(others, &usercmd.MsgPlayer{
			Id:     player.ID,
			BallId: player.SelfAnimal.GetID(),
			Name:   player.Name,
			Local:  player.GetLocation(),
			IsLive: player.IsLive,

			SnapInfo: player.GetSnapInfo(),
			Curhp:    uint32(player.SelfAnimal.GetHP()),
			Curmp:    uint32(player.SelfAnimal.GetMP()),
			Animalid: uint32(player.SelfAnimal.GetAnimalId()),
			Curexp:   player.GetExp(),
			TeamName: player.UData().TeamName,
		})
	}

	// 玩家视野中的所有球，发送给自己
	cells := scenePlayer.LookCells

	scenePlayer.LookFeeds = make(map[uint32]*bll.BallFeed)
	addfeeds, _ := scenePlayer.UpdateVeiwFeeds()
	balls = append(balls, addfeeds...)

	scenePlayer.LookBallSkill = make(map[uint32]*bll.BallSkill)
	adds, _ := scenePlayer.UpdateVeiwBallSkill()
	balls = append(balls, adds...)

	scenePlayer.LookBallFoods = make(map[uint32]*bll.BallFood)
	addfoods, _ := scenePlayer.UpdateVeiwFoods()
	balls = append(balls, addfoods...)

	//自己
	playerballs = append(playerballs, bll.PlayerBallToMsgBall(scenePlayer.SelfAnimal))
	//周围玩家
	for _, other := range scenePlayer.Others {
		if true == other.IsLive {
			playerballs = append(playerballs, bll.PlayerBallToMsgBall(other.SelfAnimal))
		}
	}

	msg := &this.Pool.MsgLoginResult
	msg.Id = scenePlayer.ID
	msg.BallId = scenePlayer.SelfAnimal.GetID()
	msg.Name = scenePlayer.Name
	msg.Ok = true
	msg.Frame = this.frame
	msg.Local = scenePlayer.GetLocation()
	msg.Balls = balls
	msg.Playerballs = playerballs
	msg.Others = others
	msg.RoomName = this.room.Name()
	msg.TeamName = scenePlayer.UData().TeamName
	msg.TeamId = scenePlayer.UData().TeamId
	msg.LeftTime = uint32(this.room.EndTime() - (time.Now().Unix() - this.room.StartTime()))
	tmp := [256]bool{0: false} //mata
	msg.BlockInf = tmp[:]
	scenePlayer.StartTime = time.Now()

	playerteam := this.room.GetTeam(scenePlayer.UData().TeamId)
	if playerteam != nil {
		msg.TeamNoticeCD = uint32(playerteam.NoticeTime - time.Now().Unix())
	}
	scenePlayer.Sess.SendCmd(usercmd.MsgTypeCmd_Login, msg)
	glog.Info("[登录] 添加玩家成功addplayer [", this.room.RoomType(), ",exp:", this.GetOfflineExp(scenePlayer.UData().Account), ",x:", x, ",y:", y, ",", this.room.ID(), ",", scenePlayer.Name, "],", scenePlayer.ID, ",", scenePlayer.UData().Account, ",ballId:", msg.BallId, ",view:", scenePlayer.GetViewRect(), ",cell:", len(cells),
		",otplayer:", len(others), "so", len(scenePlayer.Others), ",ball:", len(balls), ",teamid:", scenePlayer.UData().TeamName, ",", scenePlayer.UData().TeamId)

	this.room.BroadcastMsg(usercmd.MsgTypeCmd_AddPlayer, othermsg)

	if conf.ConfigMgr_GetMe().Global.Pystress != 0 && scenePlayer.AI.IsOK() == false {
		scenePlayer.AI.AddAICtrl()
	}

	if len(this.Players) > int(this.room.GetRoomTypeConfig().PlayerNum) {
		glog.Info("player number more than max number ", len(this.Players))
	}

	return true
}

// 删除玩家
func (this *Scene) RemovePlayer(playerId uint64) bool {
	player, ok := this.Players[playerId]
	if !ok {
		return false
	}
	oldstatus := player.IsLive
	player.OldAnimalID = player.SelfAnimal.GetAnimalId()
	this.AddOffline(player)

	this.RemoveBall(player.SelfAnimal)
	this.scenePhysic.RemoveAnimal(player.SelfAnimal.PhysicObj)

	if false == player.IsRobot {
		this.room.SaveRoomData(player)
	}
	if this.room.RoomType() == common.RoomTypeQuick || this.room.RoomType() == common.RoomTypeTeam {
		player.IsLive = false
		if this.room.RoomType() == common.RoomTypeQuick {
			player.SetDeadTime(time.Now().Unix())
			if player.IsTimeout {
				player.Dead(nil)
			}
		}
	}
	/*end*/
	if player.IsRobot {
		delete(this.Players, playerId)
	}
	glog.Info("[注销] 删除玩家成功 [timeout: ", player.IsTimeout, ",", this.room.RoomType(), ",", this.room.ID(), "],", player.ID, ",", player.UData().Account, " players:", len(this.Players), ";", oldstatus, ",exp:", player.GetExp())
	if !player.IsRobot {
		dbclient.GetMe().WriteULog("ULOG[LeaveRoom] timeout:", player.IsTimeout, ",roomType:", this.room.RoomType(), ",room:", this.room.ID(), ",id:", player.ID, ",acc:", player.UData().Account)
	}
	return true
}

func (this *Scene) GetRobotNum() int32 {
	var num int32 = 0
	for _, v := range this.Players {
		if v.IsRobot {
			num++
		}
	}
	return num
}

func (this *Scene) RemoveRobotPlayers(num int32) {
	if num > 0 {
		var ids []uint64
		for _, player := range this.Players {
			if player.IsRobot && num > 0 {
				ids = append(ids, player.ID)
				num--
			}
		}

		for _, id := range ids {
			this.room.RemovePlayerById(id)
		}
	}
}

// 团战删除机器人
func (this *Scene) DeleteTeamRoomRobot(id uint64, teamid uint32) {
	teamids := 0
	//var teamrobot []uint64
	deleteid := uint64(0)
	for _, p := range this.Players {
		if id == p.ID {
			return
		}
		if teamid == p.UData().TeamId {
			teamids = teamids + 1
			if true == p.IsRobot {
				//teamrobot = append(teamrobot, p.id)
				deleteid = p.ID
			}
		}
	}
	if teamids > 4 && 0 != deleteid {
		this.room.RemovePlayerById(deleteid)
	}
}

// 修改机器人的队伍编号 和玩家队伍名
func (this *Scene) changeTeamRobotTeamId(player *plr.ScenePlayer) {
	type teaminfo struct {
		teamname uint32
		pcount   uint8 // 真人数量
		rcount   uint8 // 机器人数量
	}
	var tmp teaminfo
	info := make(map[uint32]teaminfo)
	userteam_name := make(map[uint32]bool)
	team_name_num := 0
	for _, p := range this.Players {
		tmp.pcount = info[p.UData().TeamId].pcount
		tmp.rcount = info[p.UData().TeamId].rcount
		tmp.teamname = p.UData().TeamName
		userteam_name[p.UData().TeamName] = true
		if player.UData().TeamName == p.UData().TeamName {
			team_name_num++
		}
		if p.IsRobot {
			tmp.rcount++
		} else {
			tmp.pcount++
		}
		info[p.UData().TeamId] = tmp
	}
	// 加入已有的队伍
	val := info[player.UData().TeamId]
	if val.pcount > 1 {
		for _, p := range this.Players {
			if p.UData().TeamId == player.UData().TeamId {
				player.UData().TeamName = p.UData().TeamName
			}
			if p.IsRobot && p.UData().TeamId == player.UData().TeamId {
				//player.UData().TeamName = p.UData().TeamName
				if val.pcount+val.rcount > uint8(MAX_TEAM_USER_NUM) {
					this.room.RemovePlayerById(p.ID)
					return
				}
			}
		}
		return
	}
	// 新加的队伍没有人
	for key, tmpval := range info {
		//glog.Info("222222222 ", key, ",", tmpval.pcount, ",", tmpval.rcount, ",", tmpval.teamname)
		if tmpval.pcount == 0 {
			// 把所有机器人的队伍编号改过来  真人的队伍名改过来
			for _, p := range this.Players {
				if p.IsRobot && p.UData().TeamId == key { // 不进这里表示有问题
					if tmpval.rcount == uint8(MAX_TEAM_USER_NUM) {
						// 全是机器人需要删除一个
						tmpval.rcount--
						this.room.RemovePlayerById(p.ID)
						continue
					}
					p.UData().TeamId = player.UData().TeamId
					player.UData().TeamName = p.UData().TeamName
				}
			}
			return
		}
	}
	if val.pcount == 1 && team_name_num > 1 {
		// 自己是一个新队伍,  防止队伍名相同  如机器人新建的队伍和teamserver发来的名字相同
		player.UData().TeamName = conf.ConfigMgr_GetMe().GetTeamName(userteam_name)
		return
	}
	// 进入这边表示超过5个队伍
	if len(this.Players) > 1 {
		for key, val := range info {
			glog.Error("team_room_error ", key, ",", val.pcount, ",", val.rcount, ",", val.teamname, ",", player.ID)
		}
	}
}

// 团战加机器人
func (this *Scene) AddTeamRoomRobotNum() {
	glog.Info("AddTeamRoomRobotNum 1")
	curtime := time.Now().Unix()
	diff := curtime - this.room.StartTime()
	if curtime == diff {
		glog.Info("AddTeamRoomRobotNum 2")
		return
	}
	if diff%2 != 0 {
		glog.Info("AddTeamRoomRobotNum 3")
		return
	}
	var (
		addteam  uint32 = 0
		robotNum int    = 0
	)

	teamids := make(map[uint32]uint32)
	m_team_name := make(map[uint32]uint32)
	userteam_name := make(map[uint32]bool)
	for _, p := range this.Players {
		if p.IsRobot == true {
			robotNum = robotNum + 1
		}
		teamids[p.UData().TeamId] = teamids[p.UData().TeamId] + 1
		if 0 != p.UData().TeamName {
			m_team_name[p.UData().TeamId] = p.UData().TeamName
			userteam_name[p.UData().TeamName] = true
		}
	}
	if len(this.Players) >= MAX_TEAM_NUM*MAX_TEAM_USER_NUM {
		glog.Info("[房间补机器人] 人数已满 roomid= ", this.room.ID(), ",", len(this.Players), ",", robotNum)
		return
	}

	tmpList := common.FreeRankList{}
	for key, _ := range m_team_name {
		if teamids[key] < uint32(MAX_TEAM_USER_NUM) {
			tmpList = append(tmpList, common.FreeRank{Id: uint64(key), Score: float64(teamids[key])})
		}
	}
	// 没有5个队伍 先补上5个队伍
	if len(teamids) < MAX_TEAM_NUM {
		addteam = conf.ConfigMgr_GetMe().GetTeamName(userteam_name)
	}
	sort.Sort(tmpList)
	if 0 == addteam {
		//加入旧队伍
		teamid := len(tmpList) - 1
		if teamid < 0 {
			glog.Error("[加入旧队伍] 机器人补队伍人数不对 teamid=", teamid)
			return
		}
		glog.Info("加入旧队伍 len(tmpList)=", len(tmpList), " teamid=", teamid)
		this.AddRobotPlayer(uint32(tmpList[teamid].Id), m_team_name[uint32(tmpList[teamid].Id)])
	} else {
		//创建新的队伍
		team := len(teamids) + 1
		glog.Info("创建新的队伍 ", addteam, ",", team, ",", len(tmpList))
		this.AddRobotPlayer(uint32(team), uint32(addteam))
	}
}

func (this *Scene) CheckRobotNum() {
	if this.room.RoomType() == common.RoomTypeTeam {
		return
	}
	if conf.ConfigMgr_GetMe().Global.Pystress != 0 {
		// 压力测试，不用加机器人
		return
	}
	if conf.ConfigMgr_GetMe().Global.Norobot != 0 {
		return
	}
	rtotal := int(this.room.GetRoomTypeConfig().PlayerNum)
	rnum := len(this.Players)
	if rnum > rtotal {
		this.RemoveRobotPlayers(1)
	} else if rnum < rtotal {
		this.AddRobotPlayer(1, 1)
	}
}

func (this *Scene) SceneID() uint32 {
	return this.room.SceneID()
}

func (this *Scene) RoomType() uint32 {
	return this.room.RoomType()
}

func (this *Scene) RemoveFeed(feed *bll.BallFeed) {
	if feed.PhysicObj != nil {
		this.scenePhysic.RemoveFeed(feed.PhysicObj)
	}
}

func (this *Scene) AddMovingCube(newMovingCube *usercmd.CubeReDst) {

	this.MovingCubes[uint32(newMovingCube.CubeIndex)] = newMovingCube.RemainDistance

}

func (this *Scene) GetMovingCubes() map[uint32]int32 {
	return this.MovingCubes
}

func (this *Scene) AddFeedPhysic(feed ape.IAbstractParticle) {
	this.scenePhysic.AddFeed(feed)
}

func (this *Scene) AddAnimalPhysic(animal ape.IAbstractParticle) {
	this.scenePhysic.AddAnimal(animal)
}

func (this *Scene) RemoveAnimalPhysic(animal ape.IAbstractParticle) {
	this.scenePhysic.RemoveAnimal(animal)
}

// 地图大小（长、宽相等） XXX 改名为 SceneSize
func (this *Scene) RoomSize() float64 {
	return this.roomSize
}

func (this *Scene) UpdateSkillBallCell(ball *bll.BallSkill, oldCellID int) {
	x, y := ball.GetPos()
	newCell, ok := this.GetCell(x, y)
	if !ok || newCell.ID() == oldCellID {
		return
	}
	oldCell := this.cells[oldCellID]
	if oldCell.ID() != oldCellID {
		panic("从Cell ID获取相应Cell算法有误")
	}

	oldCell.Remove(ball.GetID(), ball.GetType())
	newCell.Add(ball)
}

func (this *Scene) Frame() uint32 {
	return this.frame
}

func (this *Scene) ResetFrame() {
	this.frame = 0
}

func (this *Scene) IncreaseFrame() {
	this.frame++
}

// 格子最大X
func (this *Scene) CellNumX() int {
	return this.cellNumX
}

// 格子最大Y
func (this *Scene) CellNumY() int {
	return this.cellNumY
}

func (this *Scene) GetPlayers() map[uint64]*plr.ScenePlayer {
	return this.Players
}
