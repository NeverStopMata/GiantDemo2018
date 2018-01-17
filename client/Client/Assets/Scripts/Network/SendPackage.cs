//****************************************************************************
// Description:发送消息对象
// Author: hiramtan@qq.com
//****************************************************************************

namespace HiGame
{
    public class SendPackage
    {
        public SendPackage(int id, byte[] body = null)
        {
            Id = id;
            Body = body;
        }

        public int Id { get; private set; }
        public byte[] Body { get; private set; }
    }
}
