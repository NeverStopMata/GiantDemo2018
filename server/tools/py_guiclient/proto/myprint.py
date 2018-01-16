import log

print=log.log


def print_login_result(pmsg):
    print("==============================")
    print("Account = ", pmsg.Account)
    print("Id = ", pmsg.Id)
    print("IsNewbie = ", pmsg.IsNewbie)
    print("Password = ", pmsg.Password)
    print("Location = ", pmsg.Location)
    print("TimeNow = ", pmsg.TimeNow)
    print("sex = ", pmsg.sex)
    print("age = ", pmsg.age)
    print("tel = ", pmsg.tel)
    print("icon = ", pmsg.icon)
    print("nick = ", pmsg.nick)
    print("channl = ", pmsg.channl)
    print("isregist = ", pmsg.isregist)
    print("PassIcon = ", pmsg.PassIcon)
    print("MaxLevel = ", pmsg.MaxLevel)
    print("idcard = ", pmsg.idcard)
    print("RegTime = ", pmsg.RegTime)
    print("")
    
    
def print_room_result(pmsg):
    print("==============================")
    print("Err = ", pmsg.Err)
    print("Addr = ", pmsg.Addr)
    print("Key = ", pmsg.Key)
    print("UId = ", pmsg.UId)
    print("Tips = ", pmsg.Tips)
    print("RoomId = ", pmsg.RoomId)
    print("Priv = ", pmsg.Priv)
    print("Model = ", pmsg.Model)
    print("SceneId = ", pmsg.SceneId)
    print("ticketnum = ", pmsg.ticketnum)
    print("")
    
    
def print_wilds_login_result(pmsg):
    print("==============================")
    print("ok = ", pmsg.ok)
    print("id = ", pmsg.id)
    print("name = ", pmsg.name)
   
def print_hit_msg(pmsg):
    print ("=============on_hit==============")
    print("source = ", pmsg.source)
    print("target = ", pmsg.target)
    print("curhp = ", pmsg.curHp)
    print("addhp = ", pmsg.addHp)
    print("")
