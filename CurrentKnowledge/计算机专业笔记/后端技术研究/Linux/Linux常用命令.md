# Linux常用命令

## 文件与目录
1. **cd命令** ：切换当前目录

> cd /home &emsp;&emsp;&emsp;&emsp;进入'/home'目录
> 
> cd ... &emsp; &emsp; &emsp; &emsp;&emsp;返回上一级目录
> 
> cd .../... &emsp; &emsp; &emsp;&emsp;返回上两级目录
> 
> cd &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;进入个人主目录
> 
> cd ~user1 &emsp;&emsp; &emsp; 进入个人的主目录
> 
> cd - &emsp; &emsp; &emsp; &emsp;  &emsp; 返回上次所在的目录

2. **pwd命令** ：显示工作路径

3. **ls命令** ：查看文件与目录的命令

> ls &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;查看目录中断文件
>
> ls -l &emsp; &emsp; &emsp; &emsp;&emsp; 显示文件和目录的详细资料
>
> ls -a &emsp; &emsp; &emsp; &emsp;&emsp;列出全部文件，包含隐藏文件
>
> ls -R &emsp; &emsp; &emsp; &emsp;&emsp;连同子目录的内容一起列出（递归列出），等于该目录下的所有文件都会显示出来
>
> ls [0-9] &emsp;&emsp; &emsp; &emsp; 显示包含数字的文件名和目录名

4. **cp命令** ：用于复制文件，copy，还可以吧多个文件一次性复制到一个目录下

> cp -a &emsp; &emsp; &emsp; &emsp;&emsp; 将文件的特性一起复制
>
> cp -p &emsp; &emsp; &emsp; &emsp;&emsp; 连同文件的属性一起复制，而非使用默认方式，与cp -a相似，常用于备份
>
> cp -i &emsp; &emsp; &emsp; &emsp;&emsp; 若目标文件已经存在，在覆盖时会询问操作的进行
>
> cp -r &emsp; &emsp; &emsp; &emsp;&emsp; 递归的持续复制，用于目录的复制行为
>
> cp -u &emsp; &emsp; &emsp; &emsp;&emsp; 目录文件与源文件有差异时才会复制

5. **mv命令** ：用于移动文件、目录或更名，move

> mv -f &emsp; &emsp; &emsp; &emsp;&emsp; force强制，如果目标文件已经存在，不会询问而直接覆盖
>
> mv -i &emsp; &emsp; &emsp; &emsp;&emsp; 若目标文件已经存在，就会询问是否覆盖
>
> mv -u &emsp; &emsp; &emsp; &emsp;&emsp; 若目标文件已经存在，且比目标文件新，才会更新

6. **rm命令** ：用于三处文件或目录，remove

> rw -f &emsp; &emsp; &emsp; &emsp;&emsp; force强制，忽略不存在的文件，不会出现警告消息
>
> rw -i &emsp; &emsp; &emsp; &emsp;&emsp; 删除文件前询问用户是否操作
>
> rw -r &emsp; &emsp; &emsp; &emsp;&emsp; 递归删除，最常用于目录删除，是一个非常危险的参数


## 查看文件内容

**cat命令** ：用于查看文本文件的内容，后接要查看的文件名，通常可用管道与more和less一起使用

> cat filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 从第一个字节开始正向查看文件内容
>
> tac filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 从最后一行开始反向查看一个文件的内容
>
> cat n filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;标示文件行数
> 
> more filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp; &emsp;查看一个长文件的内容
>
> head -n 2 filename &emsp; &emsp; &emsp; &emsp;&emsp; &emsp; &emsp; &emsp;&emsp;&emsp;&emsp; &emsp;查看一个文件的前两行
>
> tail -n 2 filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp; 查看一个文件的最后两行
> 
> tail -n +1000 filename &emsp; &emsp; &emsp; &emsp;&emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 显示文件1000行以后的内容
>
> cat filename | head -n 3000 | tail -n +1000 &emsp;&emsp; &emsp; 显示文件1000-3000行的内容
>
> cat filename | tail -n +3000 | head -n 1000 &emsp;&emsp; &emsp; 显示从3000行后共1000行的内容（3000-3999）
> 
> cat filename >> filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;将文本内容追加到其他文本中
> 
> cat -n filename > newfilename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;将文本添加行号输出一个新文件中
> 
> cat filename1 filename2 filename3 > filename4 &emsp; 合并文本内容


## 文件搜索

**find命令** ：

> find / -name filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 从'/'开始进入根文件系统搜索文件和目录
>
> find / -user username &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 搜索属于用户'username'的文件和目录
>
> find /usr/bin -type f -atime +100 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;搜索在过去100天内未被使用过的执行文件
>
> find /usr/bin -type f -mtime -10 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 搜索在十天内被创建或者修改过的文件
>
> whereis halt &emsp; &emsp; &emsp; &emsp;&emsp; &emsp; &emsp; &emsp;&emsp;&emsp;&emsp; &emsp;&emsp;&emsp;&emsp;显示一个二进制文件、源码或man的位置
>
> which halt &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;&emsp;&emsp;&emsp; 显示一个二进制文件或可执行文件的完整路径
> 
> find /var/mail/ -size +50M -exec rm {} \&emsp; &emsp; &emsp;&emsp;&emsp;删除大于50M的文件


## 文件权限

1. **chmod命令** ：使用“+”设置权限，使用“-”用于取消权限

> ls -h &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;显示权限
>
> chmod ugo+rwx directory_name &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;设置'directory_name’目录的所有人（u）、群组（g）以及其他人（o）的读（r）、写（w）和执行（x）的权限
>
> chmod go-rwx directory_name &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;删除群组（g）以及其他人（o）对'directory_name’目录的读（r）、写（w）权限

2. **chown命令** ：改变文件的所有者

> chown username filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;改变一个文件的所有人属性
>
> chown -R username directory_name &emsp; &emsp; &emsp; &emsp; &emsp;改变'directory_name’目录的所有人属性并同时改变该目录下所有文件的属性
>
> chown username:group filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 改变filename文件的所有人和群组属性

3. **chgrp命令** ：改变文件所属用户组

> chgrp group filename &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 改变filename文件的所属用户组

## 文本处理

1. **grep命令** ：分析一行的信息，若当中有我们所需要的信息（查找关键词、数字等），就将该行显示出来，该命令通常与管道命令使用，用于对一些命令的输出进行筛选加工等等

2. **paste命令** ：

> paste filename1 filename2 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 合并两个文件或两栏的内容
> 
> paste -d'+'filename1 filename2 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;合并两个文件或两栏的内容，中间用‘+’区分

3. **sort命令** ：

> sort filename1 filename2 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;排序两个文件的内容
> 
> sort filename1 filename2 | uniq &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;取出两个文件的并集
>
> sort filename1 filename2 | uniq -u &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;删除两个文件的交集，留下其他行
>
> sort filename1 filename2 | uniq -d &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;取出两个文件的交集（只留下同时存在于两个文件中的文本）

4. **comm命令** ；

> comm -1 filename1 filename2 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;比较两个文件的内容只删除filename1所包含的内容
>
> comm -2 filename1 filename2 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;比较两个文件的内容只删除filename2所包含的内容
>
> comm -3 filename1 filename2 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;比较两个文件的内容只删除两个文件共有的内容

## 打包和压缩文件

**tar命令** ：对文件进行打包，默认情况下不会压缩，如果指定了相应的参数，它还会调用相应的压缩程序（如gzip和bzip等）进行压缩和解压

> tar -c
> 
> tar -t
> 
> tar -x
> 
> tar -j
>
> tar -z
> 
> tar -v
> 
> tar -f filename
> 
> tar -C dir
> 
> tar -jcv -f filename.tar.bz2 filename
> 
> tar -jtv -f filename.tar.bz2
> 
> tar -jxc -f filename.tar.bz2 -C
> 
> bunzip2 filename.bz2
> 
> bzip2 filename
> 
> gunzip filename.gz
> 
> gzip filename
> 
> gzip -9 filename
> 
> rar a filename.rar test_file
> 
> rar a filename.rar filename1 filename2 directory_name
> 
> rar x filename.rar
> 
> zip filename.zip filename
> 
> unzip filename.zip
> 
> zip -r filename.zip filename1 filename2 directory_name

## 系统和关机

> shutdown -h now &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;关闭系统
>
> init 0 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 关闭系统
>
> telinit 0 &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;&emsp;关闭系统
> 
> shutdown -h hours:minutes & &emsp; &emsp;&emsp;按预定时间关闭系统
>
> shutdown -c &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;取消按预定时间关闭系统
>
> shutdown -r now &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 重启
> 
> reboot &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;重启
>
> logout &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;注销
>
> time &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;测算一个命令（即程序）的执行时间


## 进程相关命令

1. **jps命令** ：Java Virtual Machine Process Status Tool是JDK1.5提供的一个现实当前所有Java进程pid的命令，简单实用，用来显示当前系统的Java进程简单情况

2. **ps命令** ：Process，用于将某个时间点的进程运行情况选取下来并输出

> ps -A &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;所有进程都显示出来
>
> ps -a &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;不与terminal有关的所有进程
>
> ps -u &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;有效用户的相关进程
>
> ps -x &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;一般与a参数一起使用，可列出较完整的信息
>
> ps -l &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;较长、较详细地将PID信息列出
>
> ps aux # &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;查看系统所有的进程数据
>
> ps ax # &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp; &emsp; 查看不与terminal有关的所有进程
>
> ps -IA # &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp; 查看系统所有的进程数据
>
> ps axjf # &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp;&emsp;&emsp;查看统一部份进程树状态

3. **kill命令** ：用于向某个工作或者某个PID传送一个信号，通常与ps和jobs命令一起使用

4. **killall命令** ：向一个命令启动的进程发送一个信号

5. **top命令** ：是Linux下常用的性能分析工具，能够实时显示系统中各个进程的资源占用状况，类似与Windows的任务管理器

> 如何杀死进程：
> kill -9 PID
> killall -9 PID
> pkill PID

> 查看进程端口号：
> 
> netstat -tunlp &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; &emsp; 用于显示tcp、udp的端口和进程等相关情况
> 
> ps -ef |grep PID

> 查看端口被哪个进程占用：
> lsof -i:端口号
> 
> netstat -tunlp|grep 端口号