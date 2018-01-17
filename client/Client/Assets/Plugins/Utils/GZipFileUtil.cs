//using UnityEditor;
using System.IO;
using ICSharpCode.SharpZipLib.Zip;
using ICSharpCode.SharpZipLib.GZip;
using ICSharpCode.SharpZipLib.Checksums;
using System;
using ComponentAce.Compression.Libs.zlib;

/**
 * Author: Nash
 **/
public class GZipFileUtil
{

    /// <summary>
    /// 压缩
    /// </summary>
    public static byte[] Compress(byte[] sourceByte)
    {
        MemoryStream inStream = new MemoryStream(sourceByte);
        MemoryStream outStream = new MemoryStream();
        ZOutputStream zilbStream = new ZOutputStream(outStream, zlibConst.Z_DEFAULT_COMPRESSION);

        zilbStream.Write(sourceByte, 0, sourceByte.Length);
        zilbStream.Flush();
        zilbStream.finish();

        byte[] bytes = new byte[outStream.Length];
        outStream.Position = 0;
        outStream.Read(bytes, 0, bytes.Length);

        inStream.Close();
        outStream.Close();

        return bytes;
    }

    /// <summary>
    /// 解压缩
    /// </summary>
    public static byte[] Uncompress(byte[] sourceByte)
    {
        MemoryStream inStream = new MemoryStream(sourceByte);
        MemoryStream outStream = new MemoryStream();
        ZOutputStream zilbStream = new ZOutputStream(outStream);

        zilbStream.Write(sourceByte, 0, sourceByte.Length);
        zilbStream.Flush();
        zilbStream.finish();

        byte[] bytes = new byte[outStream.Length];
        outStream.Position = 0;
        outStream.Read(bytes, 0, bytes.Length);

        inStream.Close();
        outStream.Close();

        return bytes;
    }
    /// <summary>
    /// 使用GZIP压缩文件的方法
    /// </summary>
    /// <param name="sourcefilename">源文件路径</param>
    /// <param name="zipfilename">压缩文件路径</param>
    /// <returns>返回bool操作结果，成功true，失败 flase</returns>
    public static bool GZipFile(string sourcefilename, string zipfilename)
    {
        bool blResult;//表示压缩是否成功的返回结果
        //为源文件创建读取文件的流实例
        FileStream srcFile = File.OpenRead(sourcefilename);
        //为压缩文件创建写入文件的流实例，
        GZipOutputStream zipFile = new GZipOutputStream(File.Open(zipfilename, FileMode.Create));
        try
        {
            byte[] FileData = new byte[srcFile.Length];//创建缓冲数据
            srcFile.Read(FileData, 0, (int)srcFile.Length);//读取源文件
            zipFile.Write(FileData, 0, FileData.Length);//写入压缩文件
            blResult = true;
        }
        catch (Exception ee)
        {
            Console.WriteLine(ee.Message);
            blResult = false;
        }
        srcFile.Close();//关闭源文件
        zipFile.Close();//关闭压缩文件
        return blResult;
    }
    /// <summary>
    /// 使用GZIP解压文件的方法
    /// </summary>
    /// <param name="zipfilename">源文件路径</param>
    /// <param name="unzipfilename">解压缩文件路径</param>
    /// <returns>返回bool操作结果，成功true，失败 flase</returns>
    public static bool UnGzipFile(string zipfilename, string unzipfilename)
    {
        bool blResult;//表示解压是否成功的返回结果
        //创建压缩文件的输入流实例
        GZipInputStream zipFile = new GZipInputStream(File.OpenRead(zipfilename));
        //创建目标文件的流
        FileStream destFile = File.Open(unzipfilename, FileMode.Create);
        try
        {
            int buffersize = 2048;//缓冲区的尺寸，一般是2048的倍数
            byte[] FileData = new byte[buffersize];//创建缓冲数据
            while (buffersize > 0)//一直读取到文件末尾
            {
                buffersize = zipFile.Read(FileData, 0, buffersize);//读取压缩文件数据
                destFile.Write(FileData, 0, buffersize);//写入目标文件
            }
            blResult = true;
        }
        catch (Exception ee)
        {
            Console.WriteLine(ee.Message);
            blResult = false;
        }
        destFile.Close();//关闭目标文件
        zipFile.Close();//关闭压缩文件
        return blResult;
    }

    /// <summary>
    /// 压缩单个文件
    /// </summary>
    /// <param name="fileToZip">要压缩的文件</param>
    /// <param name="zipedFile">压缩后的文件</param>
    /// <param name="compressionLevel">压缩等级</param>
    /// <param name="blockSize">每次写入大小</param>
    public static void ZipFile(string fileToZip, string zipedFile, int compressionLevel, int blockSize)
    {
        //如果文件没有找到，则报错
        if (!System.IO.File.Exists(fileToZip))
        {
            throw new System.IO.FileNotFoundException("指定要压缩的文件: " + fileToZip + " 不存在!");
        }

        using (System.IO.FileStream ZipFile = System.IO.File.Create(zipedFile))
        {
            using (ZipOutputStream ZipStream = new ZipOutputStream(ZipFile))
            {
                using (System.IO.FileStream StreamToZip = new System.IO.FileStream(fileToZip, System.IO.FileMode.Open, System.IO.FileAccess.Read))
                {
                    string fileName = fileToZip.Substring(fileToZip.LastIndexOf("\\") + 1);

                    ZipEntry ZipEntry = new ZipEntry(fileName);

                    ZipStream.PutNextEntry(ZipEntry);

                    ZipStream.SetLevel(compressionLevel);

                    byte[] buffer = new byte[blockSize];

                    int sizeRead = 0;

                    try
                    {
                        do
                        {
                            sizeRead = StreamToZip.Read(buffer, 0, buffer.Length);
                            ZipStream.Write(buffer, 0, sizeRead);
                        }
                        while (sizeRead > 0);
                    }
                    catch (System.Exception ex)
                    {
                        throw ex;
                    }

                    StreamToZip.Close();
                }

                ZipStream.Finish();
                ZipStream.Close();
            }

            ZipFile.Close();
        }
    }

    /// <summary>
    /// 压缩单个文件
    /// </summary>
    /// <param name="fileToZip">要进行压缩的文件名</param>
    /// <param name="zipedFile">压缩后生成的压缩文件名</param>
    public static void ZipFile(string fileToZip, string zipedFile)
    {
        //如果文件没有找到，则报错
        if (!File.Exists(fileToZip))
        {
            throw new System.IO.FileNotFoundException("指定要压缩的文件: " + fileToZip + " 不存在!");
        }

        using (FileStream fs = File.OpenRead(fileToZip))
        {
            byte[] buffer = new byte[fs.Length];
            fs.Read(buffer, 0, buffer.Length);
            fs.Close();

            using (FileStream ZipFile = File.Create(zipedFile))
            {
                using (ZipOutputStream ZipStream = new ZipOutputStream(ZipFile))
                {
                    string fileName = fileToZip.Substring(fileToZip.LastIndexOf("\\") + 1);
                    ZipEntry ZipEntry = new ZipEntry(fileName);
                    ZipStream.PutNextEntry(ZipEntry);
                    ZipStream.SetLevel(5);

                    ZipStream.Write(buffer, 0, buffer.Length);
                    ZipStream.Finish();
                    ZipStream.Close();
                }
            }
        }
    }

    /// <summary>
    /// 压缩多层目录
    /// </summary>
    /// <param name="strDirectory">The directory.</param>
    /// <param name="zipedFile">The ziped file.</param>
    public static void ZipFileDirectory(string strDirectory, string zipedFile)
    {
        using (System.IO.FileStream ZipFile = System.IO.File.Create(zipedFile))
        {
            using (ZipOutputStream s = new ZipOutputStream(ZipFile))
            {
                ZipSetp(strDirectory, s, "");
            }
        }
    }

    /// <summary>
    /// 递归遍历目录
    /// </summary>
    /// <param name="strDirectory">The directory.</param>
    /// <param name="s">The ZipOutputStream Object.</param>
    /// <param name="parentPath">The parent path.</param>
    private static void ZipSetp(string strDirectory, ZipOutputStream s, string parentPath)
    {
        //if (strDirectory[strDirectory.Length - 1] != Path.DirectorySeparatorChar)
        //{
        //    strDirectory += Path.DirectorySeparatorChar;
        //}
        //Crc32 crc = new Crc32();

        //string[] filenames = Directory.GetFileSystemEntries(strDirectory);

        //foreach (string file in filenames)// 遍历所有的文件和目录
        //{

        //    if (Directory.Exists(file))// 先当作目录处理如果存在这个目录就递归Copy该目录下面的文件
        //    {
        //        string pPath = parentPath;
        //        pPath += file.Substring(file.LastIndexOf("\\") + 1);
        //        pPath += "\\";
        //        ZipSetp(file, s, pPath);
        //    }

        //    else // 否则直接压缩文件
        //    {
        //        //打开压缩文件
        //        using (FileStream fs = File.OpenRead(file))
        //        {
        //            byte[] buffer = new byte[fs.Length];
        //            fs.Read(buffer, 0, buffer.Length);

        //            string fileName = parentPath + file.Substring(file.LastIndexOf("\\") + 1);
        //            ZipEntry entry = new ZipEntry(fileName);

        //            entry.DateTime = DateTime.Now;
        //            entry.Size = fs.Length;

        //            fs.Close();

        //            crc.Reset();
        //            crc.Update(buffer);

        //            entry.Crc = crc.Value;
        //            s.PutNextEntry(entry);

        //            s.Write(buffer, 0, buffer.Length);
        //        }
        //    }
        //}
    }

    /// <summary>
    /// 解压缩一个 zip 文件。
    /// </summary>
    /// <param name="zipedFile">The ziped file.</param>
    /// <param name="strDirectory">The STR directory.</param>
    /// <param name="password">zip 文件的密码。</param>
    /// <param name="overWrite">是否覆盖已存在的文件。</param>
    public static void UnZip(string zipedFile, string strDirectory, string password = "0000", bool overWrite = true)
    {
        if (strDirectory == "")
        {
            strDirectory = Directory.GetCurrentDirectory();
        }   
        using (ZipInputStream s = new ZipInputStream(File.OpenRead(zipedFile)))
        {
            s.Password = password;
            ZipEntry theEntry;
            while ((theEntry = s.GetNextEntry()) != null)
            {
                string directoryName = "";
                string pathToZip = "";
                pathToZip = theEntry.Name;
                if (pathToZip != "")
                {
                    directoryName = Path.GetDirectoryName(pathToZip);
                }
                string fileName = Path.GetFileName(pathToZip);
                string targetDirectoryUrl = System.IO.Path.Combine(strDirectory, directoryName);
                Directory.CreateDirectory(targetDirectoryUrl);
                if (fileName != "")
                {
                    string targetFileUrl = System.IO.Path.Combine(targetDirectoryUrl, fileName);
                    if ((File.Exists(targetFileUrl) && overWrite) || (!File.Exists(targetFileUrl)))
                    {
                        using (FileStream streamWriter = File.Create(targetFileUrl))
                        {
                            int size = 2048;
                            byte[] data = new byte[2048];
                            while (true)
                            {
                                size = s.Read(data, 0, data.Length);
                                if (size > 0)
                                    streamWriter.Write(data, 0, size);
                                else
                                    break;
                            }
                            streamWriter.Close();
                        }
                    }
                }
            }
            s.Close();
        }
    }
}
