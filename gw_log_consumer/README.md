# 日志自定义解析器

之前用python实现过一个版本，拿golang练练手。

## 用途

解析已经从filebeat送入到redis的特定格式数据：

```
2017-10-11 12:13:14 121231231 ELK_DATA {"unixtimestamp":15,"@timestamp":"","submitrecv":100,"submitallow":55} ELK_END.
```

然后将数据写出到logstash监听的unixsock，最终由logstash送入elasticsearch.


