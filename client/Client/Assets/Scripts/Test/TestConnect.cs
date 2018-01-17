using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using HiGame;
using usercmd;
public class TestConnect : MonoBehaviour {

    MsgDoNothing msg;
    TcpConnection conn;
    // Use this for initialization
    void Start()
    {
         msg = new MsgDoNothing
        {
            hello = "hello world",
            id = 9
        };
        Debuger.EnableOnConsole(true);
        conn = TcpConnection.Instance;
        conn.ServiceEventHandler = ServiceOk;
    }


    void ServiceOk()
    {
        conn.Connect("192.168.251.66", 9001);
        Debuger.Log("connect");

    }
    // Update is called once per frame
    void Update()
    {
        if(Input.GetMouseButtonDown(0))
            MsgHandler.Send((int)MsgTypeCmd.DoNothing, msg);
    }
}
