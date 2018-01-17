//****************************************************************************
// Description:消息注册发送及处理
// Author: hiramtan@qq.com
//****************************************************************************
using GameBox.Framework;
using System;
using System.Collections.Generic;

namespace HiGame
{
    public class MsgHandler
    {
        private static Dictionary<int, Action<IProto>> _dic = new Dictionary<int, Action<IProto>>();

        public static void Regist(int id, Action<IProto> action)
        {
            if (_dic.ContainsKey((int)id))
            {
                //Debuger.LogWarning("dont need regist again");
                return;
            }
            _dic.Add((int)id, action);
        }

        public static void UnRegist(int id)
        {
            if (_dic.ContainsKey((int)id))
            {
                _dic.Remove((int)id);
                return;
            }
            AnyLogger.E("do not have key:" + id);
        }

        public static void Dispatch(int id, byte[] bytes)
        {
            if (!_dic.ContainsKey(id))
            {
                //Debuger.LogWarning("should regeist first:" + id);
                return;
            }
            IProto iProto = new Proto(bytes);
            _dic[id](iProto);
        }

        private static void SendBytes(int id, byte[] bytes = null)
        {
            SendPackage sendPackage = bytes == null ? new SendPackage((int)id) : new SendPackage((int)id, bytes);
            if (TcpConnection.Instance == null)
            {
                AnyLogger.E("should init socket first");
                return;
            }
            TcpConnection.Instance.Send(sendPackage);
        }

        public static void Send(int id, object obj = null)
        {
            if (obj == null)
            {
                SendBytes(id);
            }
            else
            {
                IProto t = new Proto(obj);
                SendBytes(id, t.Get());
            }
        }

   

        public static void Clear()
        {
            _dic.Clear();
        }
    }
}