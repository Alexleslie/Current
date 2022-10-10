# 协议升级机制

HTTP 协议 提供了一种特殊的机制，这一机制允许将一个已建立的连接升级成新的、不相容的协议。这篇指南涵盖了其工作原理和使用场景。

通常来说这一机制总是由客户端发起的（不过也有例外，比如说可以由服务端发起升级到传输层安全协议TLS），服务端可以选择是否要升级到新协议。借助这一技术，连接可以以常用的协议启动（如 HTTP/1.1），随后再升级到 HTTP2 甚至是 WebSockets.

> HTTP/2 明确禁止使用此机制，这个机制只属于 HTTP/1.1


## 协议升级
协议升级请求总是由客户端发起的；暂时没有服务端请求协议更改的机制。当客户端试图升级到一个新的协议时，可以先发送一个普通的请求（GET，POST等），不过这个请求需要进行特殊配置以包含升级请求。

这个请求需要添加两项额外的 header：
- Connection: Upgrade, 设置 Connection 头的值为 "Upgrade" 来指示这是一个升级请求。 
- Upgrade: protocols, Upgrade 头指定一项或多项协议名，按优先级排序，以逗号分隔。

一个典型的包含升级请求的例子差不多是这样的：

    GET /index.html HTTP/1.1
    Host: www.example.com
    Connection: upgrade
    Upgrade: example/1, foo/2

服务在发送 101 状态码之后，就可以使用新的协议，并可以根据需要执行任何其他协议指定的握手。实际上，一旦这次升级完成了，连接就变成了双向管道。并且可以通过新协议完成启动升级的请求。

>  HTTP/2 已经不再支持 101 状态码了，也不再支持任何连接升级机制

## 升级机制的常用场合

### HTTP升级为HTTPS
#### 客户端
客户端可以主动将 HTTP/1.1 连接升级到 TLS/1.0。这样做的主要优点是可以避免在服务器上使用从“http://”到“https://”的 URL 重定向，并且可以轻松地在虚拟主机上使用 TLS。但是，这可能会给代理服务器带来问题。

    GET http://destination.server.ext/secretpage.html HTTP/1.1
    Host: destination.server.ext
    Upgrade: TLS/1.0
    Connection: Upgrade

如果服务器确实支持 TLS 升级并希望允许升级，它会使用 101 Switching Protocols 响应代码进行响应，如下所示：

    HTTP/1.1 101 Switching Protocols
    Upgrade: TLS/1.0, HTTP/1.1

对 TLS 的请求可以是可选的，也可以是强制的。

#### 服务端主动升级
服务器启动升级到 TLS
这与客户端启动的升级方式大致相同；通过添加Upgrade到标头来请求可选升级。但是，强制升级的工作方式略有不同，因为它通过回复收到的带有426状态码的消息来请求升级，如下所示：

    HTTP/1.1 426 Upgrade Required
    Upgrade: TLS/1.1, HTTP/1.1
    Connection: Upgrade
    
    <html>
    ... Human-readable HTML page describing why the upgrade is required
        and what to do if this text is seen ...
    </html>

如果接收到 426 Upgrade Required 响应的客户端愿意并且能够升级到 TLS，那么它应该启动客户端发起 TLS 升级的相同过程。