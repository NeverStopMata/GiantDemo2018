//****************************************************************************
// Description:proto序列反序列化
// Author: hiramtan@qq.com
//****************************************************************************
using System.IO;
using ProtoBuf;

namespace HiGame
{
    public class Proto : IProto
    {
        private readonly byte[] _bytes;
        private readonly object _obj;

        public Proto(object obj)
        {
            _obj = obj;
        }

        public Proto(byte[] bytes)
        {
            _bytes = bytes;
        }

        public T Get<T>()
        {
            return Deserialize<T>(_bytes);
        }

        public byte[] Get()
        {
            return Serialize(_obj);
        }

        private byte[] Serialize<T>(T obj)
        {
            using (var steam = new MemoryStream())
            {
                Serializer.Serialize(steam, obj);
                return steam.ToArray();
            }
        }

        private T Deserialize<T>(byte[] bytes)
        {
            using (var stream = new MemoryStream(bytes))
            {
                var obj = default(T);
                obj = Serializer.Deserialize<T>(stream);
                return obj;
            }
        }
    }

    public interface IProto
    {
        T Get<T>();

        byte[] Get();
    }
}