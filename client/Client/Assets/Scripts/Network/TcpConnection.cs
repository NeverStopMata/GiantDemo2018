//****************************************************************************
// Description:tcp通信逻辑
// Author: hiramtan@qq.com
//****************************************************************************
using GameBox.Framework;
using GameBox.Service.ByteStorage;
using GameBox.Service.NetworkManager;
using System;
using usercmd;
using UnityEngine;
namespace HiGame
{
    public class TcpConnection : Connection, IProtocolHost
    {
        public static TcpConnection Instance = new TcpConnection();


        public Action ServiceEventHandler;


        public enum EMessageKey
        {
            //TeamServer,
            RoomServer
        }
        private const string Connected = "connected";
        private const string Connecting = "connecting";
        private const string Disconnected = "disconnected";
        private INetworkManager _iNetworkManager;
        private ISocketChannel _iSocketChannel;
        private bool _isRoomServerLogin;
        private bool _isTeamServerLogin;

        //private HeartBeatMsg heartBeatMsg;
        public TcpConnection()
        {
            //GameWorld.Instance.OnApplicationQuit(Disconnect);
            new ServiceTask(new[]
            {
                typeof(IByteStorage),
                typeof(INetworkManager)
            }).Start().Continue(t =>
            {
                _iNetworkManager = ServiceCenter.GetService<INetworkManager>();


                _iSocketChannel = _iNetworkManager.Create("tcp") as ISocketChannel;
                _iSocketChannel.Setup(this);
                _iSocketChannel.OnChannelStateChange = OnServerStateChange;

                if (ServiceEventHandler!=null)
                    ServiceEventHandler();
              //  MsgHandler.Send((int)usercmd.MsgTypeCmd.Login, new usercmd.MsgLogin());
                return null;
            });
        }

        public void Pack(IObjectReader reader, IByteArray writer)
        {
            var message = reader.ReadOne() as SendPackage;
            var byteArray = Pack(message);
            writer.WriteBytes(byteArray.Bytes, 0, byteArray.Length);
        }

        public void Unpack(IByteArray reader, IObjectWriter writer)
        {
            _receiveArray.WriteBytes(reader.ReadBytes());
            Unpack();
        }

        public void Connect(string ip, int port)
        {
            _iSocketChannel.Connect(ip, port);
        }

        public void Disconnect()
        {
            //if (heartBeatMsg != null)
            //{
            //    heartBeatMsg.OnClose();
            //}
            //if (_iSocketChannel != null)
            //{
            //    _iSocketChannel.Disconnect();
            //    _iSocketChannel.Dispose();
            //    _iSocketChannel = null;
            //}
        }

        public void Send(SendPackage obj)
        {
            _iSocketChannel.Send(obj);
        }

        private void OnServerStateChange(string state)
        {
            switch (state)
            {
                case Connected:
                    Debuger.Log("OnTeamServerState" + Connected);
                    //ServerConnected();
                    break;
                case Connecting:
                    Debuger.Log("OnTeamServerState" + Connecting);
                    break;
                case Disconnected:
                    //var str = Singleton.GetInstance<Localization>().Get(ExcelManager.settings_error[1].error);
                    //Service.UIManager.ShowPop(str);
                    Debuger.Log("OnTeamServerState" + Disconnected);
                    break;
            }
        }

        private void ServerConnected()
        {
            //heartBeatMsg = new HeartBeatMsg((int)MsgTypeCmd.HeartBeat);
            //new RoomLoginMsg((int)MsgTypeCmd.Login);
        }
    }
}