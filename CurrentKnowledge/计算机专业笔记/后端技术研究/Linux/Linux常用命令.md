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
> cp -p &emsp; &emsp; &emsp; &emsp;&emsp; 连同问津的属性一起复制，而非使用默认方式，与cp -a相似，常用于备份
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
> rw -i &emsp; &emsp; &emsp; &emsp;&emsp; 删除问津前询问用户是否操作
>
> rw -r &emsp; &emsp; &emsp; &emsp;&emsp; 递归删除，最常用于目录删除，是一个非常危险的参数


## 查看文件内容

**cat命令** ：

## 文件搜索

**find命令** ：

## 文件权限

使用“+”设置权限，使用“-”用于取消权限
1. **chmod命令** ：
2. **chown命令** ：
3. **chgrp命令** ：

## 文本处理

1. **grep命令** ：
2. **paste命令** ：
3. **sort命令** ：
4. **comm命令** ；

## 打包和压缩文件

**tar命令** ：

## 系统和关机


## 进程相关命令

1. **jps命令** ：
2. **ps命令** ：
3. **kill命令** ：
4. **killall命令** ：
5. **top命令** ：