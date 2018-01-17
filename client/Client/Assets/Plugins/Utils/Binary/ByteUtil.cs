public static class ByteUtil
{
    public static bool GetBit(uint value, int index)
    {
        return value == (value | (uint)(1 << index));
    }

    public static uint SetBit(uint value, int index, bool boolean)
    {
        if(boolean)
        {
            value = (value | (uint)(1 << index));
        }
        else
        {
            value = (value & ~(uint)(1 << index));
        }
        return value;
    }

    public static bool[] ToBitArray(uint value, int length)
    {
        bool[] temp = new bool[length];
        for(int i= 0; i< length; i++)
        {
            temp[i] = (value & 1) == 1;
            value = value >> 1;
        }
        return temp;
    }
}
