using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using UnityEngine;

public class ByteArray
{
    private const int BYTE_LEN = 1;
    private const int USHORT_LEN = 2;
    private const int UINT_LEN = 4;
    private const int LONG_LEN = 8;
    private const int FLOAT_LEN = 4;
    private const int INT_LEN = 4;

    private int _length = 0;
    private List<byte> _byteList;

    public ByteArray()
    {
        _byteList = new List<byte>();
    }

    public void WriteBoolean(bool value)
    {
        _length = _length + BYTE_LEN;
        _byteList.Add(Convert.ToByte(value));
    }

    public void WriteUnsignedByte(uint value)
    {
        _length = _length + BYTE_LEN;
        _byteList.Add(Convert.ToByte(value));
    }

    public void WriteUnsignedShort(uint value)
    {
        _length = _length + USHORT_LEN;
        byte[] temp = BitConverter.GetBytes(value);

        for(int i= 0; i< USHORT_LEN; i++)
        {
            _byteList.Add(temp[i]);
        }
    }

    public void WriteUnsignedInt(uint value)
    {
        _length = _length + UINT_LEN;
        byte[] temp = BitConverter.GetBytes(value);
        
        for(int i= 0; i< UINT_LEN; i++)
        {
            _byteList.Add(temp[i]);
        }
    }

    public void WriteSignedInt(int value)
    {
        _length = _length + INT_LEN;
        byte[] temp = BitConverter.GetBytes(value);

        for (int i = 0; i < INT_LEN; i++)
        {
            _byteList.Add(temp[i]);
        }
    }

    public void WriteUnsignedLong(ulong value)
    {
        _length = _length + LONG_LEN;
        byte[] temp = BitConverter.GetBytes(value);
        for(int i = 0; i< LONG_LEN; i++)
        {
            _byteList.Add(temp[i]);
        }
    }

    public void WriteUTFBytes(string value, int length)
    {
        _length = _length + length;
        byte[] temp = Encoding.UTF8.GetBytes(value);

        for(int i = 0; i< length; i++)
        {
            if(i < temp.Length)
            {
                _byteList.Add(temp[i]);
            }
            else
            {
                _byteList.Add(0);
            }
        }
    }

    public void WriteBytes(byte[] value, int length = 0)
    {
        if(length > 0)
        {
            _length = _length + length;

            for(int i = 0; i< length; i++)
            {
                if(i < value.Length)
                {
                    _byteList.Add(value[i]);
                }
                else
                {
                    _byteList.Add(0);
                }
            }
        }
        else
        {
            _length = _length + value.Length;
            _byteList.AddRange(value);
        }
    }

    public void WriteBytes(ByteArray value, int length = 0)
    {
        if(length > 0)
        {
            _length = _length + length;

            for(int i = 0; i< length; i++)
            {
                if(i < value._length)
                {
                    _byteList.Add(value._byteList[i]);
                }
                else
                {
                    _byteList.Add(0);
                }
            }
            value._length = value._length - length;
            value._byteList.RemoveRange(0, length);
        }
        else
        {
            _length = _length + value._length;
            _byteList.AddRange(value._byteList);
        }
    }

    public bool ReadBoolean()
    {
        if(_length < BYTE_LEN)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return false;
        }
        else
        {
            byte temp = _byteList[0];
            _length = _length - BYTE_LEN;
            _byteList.RemoveRange(0, BYTE_LEN);
            return Convert.ToBoolean(temp);
        }
    }

    public uint ReadUnsignedByte()
    {
        if (_length < BYTE_LEN)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return 0;
        }
        else
        {
            byte temp = _byteList[0];
            _length = _length - BYTE_LEN;
            _byteList.RemoveRange(0, BYTE_LEN);
            return Convert.ToUInt32(temp);
        }
    }

    public uint ReadUnsignedShort()
    {
        if(_length < USHORT_LEN)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return 0;
        }
        else
        {
            byte[] temp = new byte[UINT_LEN];

            for(int i=0; i< UINT_LEN; i++)
            {
                if(i < USHORT_LEN)
                {
                    temp[i] = _byteList[i];
                }
                else
                {
                    temp[i] = 0;
                }
            }
            _length = _length - USHORT_LEN;
            _byteList.RemoveRange(0, USHORT_LEN);

            return BitConverter.ToUInt32(temp, 0);
        }
    }

    public uint ReadUnsignedInt()
    {
        if(_length < UINT_LEN)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return 0;
        }
        else
        {
            byte[] temp = new byte[UINT_LEN];

            for (int i = 0; i < UINT_LEN; i++)
            {
                temp[i] = _byteList[i];
            }
            _length = _length - UINT_LEN;
            _byteList.RemoveRange(0, UINT_LEN);

            return BitConverter.ToUInt32(temp, 0);
        }
    }
    public int ReadSignedInt()
    {
        if (_length < INT_LEN)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return 0;
        }
        else
        {
            byte[] temp = new byte[INT_LEN];

            for (int i = 0; i < INT_LEN; i++)
            {
                temp[i] = _byteList[i];
            }
            _length = _length - INT_LEN;
            _byteList.RemoveRange(0, INT_LEN);

            return BitConverter.ToInt32(temp, 0);
        }
    }

    public float ReadFloat()
    {
        if(_length < FLOAT_LEN)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return 0f;
        }
        else
        {
            byte[] temp = new byte[FLOAT_LEN];

            for(int i=0; i< FLOAT_LEN; i++)
            {
                temp[i] = _byteList[i];
            }

            _length = _length - FLOAT_LEN;
            _byteList.RemoveRange(0, FLOAT_LEN);
            return BitConverter.ToSingle(temp, 0);
        }
    }

    public ulong ReadUnsignedLong()
    {
        if(_length < LONG_LEN)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return 0;
        }
        else
        {
            byte[] temp = new byte[LONG_LEN];
            for(int i= 0; i< LONG_LEN; i++)
            {
                temp[i] = _byteList[i];
            }
            _length = _length - LONG_LEN;
            _byteList.RemoveRange(0, LONG_LEN);
            return BitConverter.ToUInt64(temp, 0);
        }
    }

    public string ReadUTFByte(int length)
    {
        if (_length < length)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
            return "";
        }
        else
        {
            byte[] temp = new byte[length];
            for(int i= 0; i< length; i++)
            {
                temp[i] = _byteList[i];
            }

            _length = _length - length;
            _byteList.RemoveRange(0, length);

            int removeCount = 0;
            for(int j = temp.Length - 1; j >= 0; j--)
            {
                if(temp[j] == 0)
                {
                    removeCount++;
                }
                else
                {
                    break;
                }
            }
            return Encoding.UTF8.GetString(temp, 0, temp.Length - removeCount);
        }
    }

    public void ReadBytes(byte[] value, int length = 0)
    {
        if(_length < length)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
        }
        else
        {
            if(length > 0)
            {
                for(int i=0; i< length; i++)
                {
                    value[i] = _byteList[i];
                }
                _length = _length - length;
                _byteList.RemoveRange(0, length);
            }
            else
            {
                _byteList.CopyTo(value);
                Clear();
            }
        }
    }

    public void ReadBytes(ByteArray value, int length = 0)
    {
        if(_length < length)
        {
            Debug.LogException(new Exception("ByteArray == Not Enough Data Length"));
        }
        else
        {
            if(length > 0)
            {
                value._length = value._length + length;
                for(int i= 0; i< length; i++)
                {
                    value._byteList.Add(_byteList[i]);
                }
                _length = _length - length;
                _byteList.RemoveRange(0, length);
            }
            else
            {
                value._length = value._length + _length;
                value._byteList.AddRange(_byteList);
                Clear();
            }
        }
    }

    public MemoryStream Stream
    {
        get
        {
            MemoryStream stream = new MemoryStream();
            stream.Write(_byteList.ToArray(), 0, _length);
            stream.Position = 0;
            return stream;
        }
    }

    public byte[] Bytes
    {
        get
        {
            return _byteList.ToArray();
        }
    }

    public int Length
    {
        get
        {
            return _length;
        }

        set
        {
            _length = value;
            if(_byteList.Count > _length)
            {
                _byteList.RemoveRange(_length, _byteList.Count - _length);
            }
            else if (_byteList.Count < _length)
            {
                while(_byteList.Count <  _length)
                {
                    _byteList.Add(0);
                }
            }
        }
    }

    public void Compress()
    {
        byte[] temp = GZipFileUtil.Compress(Bytes);

        Clear();
        WriteBytes(temp);
    }

    public void UnCompress()
    {
        byte[] temp = GZipFileUtil.Uncompress(Bytes);

        Clear();
        WriteBytes(temp);
    }

    public void Clear()
    {
        _length = 0;
        _byteList.Clear();
    }

    public ByteArray Clone()
    {
        ByteArray temp = new ByteArray();
        temp._length = _length;
        temp._byteList.AddRange(_byteList);
        return temp;
    }

    public static ByteArray Get()
    {
        if (_cache.Count > 0)
            return _cache.Dequeue();
        return new ByteArray();
    }

    public static void Cache(ByteArray byteArr)
    {
        byteArr.Clear();
        _cache.Enqueue(byteArr);
    }

    private static Queue<ByteArray> _cache = new Queue<ByteArray>();
}
