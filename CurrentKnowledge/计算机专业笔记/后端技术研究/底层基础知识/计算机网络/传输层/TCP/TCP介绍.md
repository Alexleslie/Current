# TCP介绍
TCP是在**不可靠**的IP层上实现可靠的数据传输协议。主要解决传输的可靠，有序，无丢失，不重复的问题

## 特点
1. TCP是面向连接的传输层协议
2. TCP连接只能是点对点的（一个进程对应一个进程）
3. TCP提供可靠的交付，保证数据无差错，无丢失，有序
4. TCP提供全双工通信，允许双方在任何时候都可以发送数据
5. TCP是面向字节流的，TCP把应用程序的数据仅仅视为一连串无结构的字符串

## 优点
1. 可靠的数据传输
2. 网络的拥塞控制
3. 按序发送数据

## 缺点
1. 连接建立比较慢 
2. 慢启动，刚开始传输数据时速率很慢 
3. 拥塞会降低传输速率
4. 数据流阻塞，假设发送的数据流有一段不可用，即使后面的数据可用，但也不能提前发送
