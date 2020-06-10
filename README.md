### 组件

nsq由三大守护进程组成

- **nsqd**   负责接受、排列和分发消息
- **nsqlookup**  提供nsqd信息
- **Nsqadmin**  监控和分析消息队列集群

### 基本架构

nsq集群基本架构图

![](http://image.mingjing.xyz/blog/nsq.png)

单个nsqd工作原理：图片来源官网

![](https://f.cloud.github.com/assets/187441/1700696/f1434dc8-6029-11e3-8a66-18ca4ea10aca.gif)

- producer通过tcp/http方式，将生产的消息发送给指定nsqd peer的topic，若topic不存在则创建。topic又备份消息一份各给它的所有 channel 。
- tcp client 从channel读取消息并发送给 consumer 消费，若一个channel被多个client订阅，则消息谁抢到则谁消费
- consumer消费完通知nsqd peer；若消费失败可要求重发。所以nsq的消息只能保证一定消费一次。

**集群nsqd**

集群nsqd中每个nsqd也是独立的，它只是把信息注册到nsqlookup，这样consumer可以到nsqlookup查找nsqd信息，进行连接。比如consumer1想消费<topic1,channel1>的消息，它就到nsqlookup查找所有拥有该channel的nsqd peer 的信息，然后通过tcp连接nsqd peer。

### 原理分析



### 总结







## 参考

[nsq官方文档](https://nsq.io/overview/internals.html)

https://zhuanlan.zhihu.com/p/37081073

