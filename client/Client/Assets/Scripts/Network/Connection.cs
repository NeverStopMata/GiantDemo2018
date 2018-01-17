
namespace HiGame
{
    public abstract class Connection
    {
        public enum ESocket
        {
            Tcp,
            Kcp,
            Udp
        }
        // public static string IP = "122.11.58.160:8123";
        public static string IP = "127.0.0.1:8123";
        public static ESocket ESocketType = ESocket.Udp;
        protected const int MaxPackageLength = 65535;
        protected const int PackageHeadLength = 6;
        protected const int CommandHeadLength = 2;
        protected readonly ByteArray _bodyArray = new ByteArray();
        protected readonly ByteArray _receiveArray = new ByteArray();
        protected int cmdId;
        protected bool isCompress;
        protected bool isHaveHead;
        protected uint packBodyLen;

        public Connection()
        {
           // Ping();
        }

        protected ByteArray Pack(SendPackage obj)
        {
            return null;
            //return CmdUtil.GetPackage(obj.Id, obj.Body);
        }

        protected void Unpack()
        {
            while (true)
                if (!isHaveHead)
                {
                    //三个字节 数据块大小  一个字节0 消息全部未压缩   两个字节消息id
                    if (_receiveArray.Length > PackageHeadLength)
                    {
                        var len1 = _receiveArray.ReadUnsignedByte();
                        var len2 = _receiveArray.ReadUnsignedByte();
                        var len3 = _receiveArray.ReadUnsignedByte();
                        packBodyLen = len1 | (len2 << 8) | (len3 << 16);
                        isCompress = _receiveArray.ReadUnsignedByte() == 1;
                        cmdId = (int)_receiveArray.ReadUnsignedShort();
                        if (packBodyLen < CommandHeadLength || packBodyLen > MaxPackageLength)
                        {
                            _receiveArray.Clear();
                            break;
                        }
                        if (packBodyLen == CommandHeadLength)
                            MsgHandler.Dispatch(cmdId, null);
                        else
                            isHaveHead = true;
                    }
                    else
                    {
                        break;
                    }
                }
                else
                {
                    if (_receiveArray.Length >= packBodyLen - CommandHeadLength)
                    {
                        _bodyArray.Clear();
                        _receiveArray.ReadBytes(_bodyArray, (int)packBodyLen - CommandHeadLength);
                        if (isCompress)
                            _bodyArray.UnCompress();
                        MsgHandler.Dispatch(cmdId, _bodyArray.Bytes);
                        isHaveHead = false;
                    }
                    else
                    {
                        break;
                    }
                }
        }

        public virtual void DisConnect()
        {

        }
#if HExpectation
        public static uint Frame { get; private set; }
        private static int time;
        public static void SetFrame(uint frame)
        {
            Frame = frame;
            time = Mathf.RoundToInt(Time.realtimeSinceStartup * 1000);
        }

        public static uint GetExpectation()
        {
            if (Frame == 0)
                return 1;
            int now = Mathf.RoundToInt(Time.realtimeSinceStartup * 1000);
            var value = Frame + Mathf.FloorToInt((now - time) / 100) * 4;
            return (uint)value;
        }
#endif
        //public static int pingTime;
        //private static void Ping()
        //{
        //    var ip = IP.Split(':')[0];
        //    new AsyncPingTask((x) => { pingTime = x; }, ip, Common.pingTime).Start();
        //}
    }
}