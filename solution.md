# tinykv课程思路
## project 1
1. `kv/storage/standalone_storage/standalone_storage.go`
    * 设计`StanAloneStorage`的结构：使用`package engine_util`提供的`Engines`类型即可。
    * `Reader`和`Write`方法实现：

### 疑问
1. 参照`kv/storage/raft_storage/region_reader.go`中写法，`badger.Txn.Discard()`方法只在`Reader`的`Close()`方法中被调用，而测试代码`server_test.go`中并没有调用`Close()`去提交事务，测试代码这样写虽然没问题，但是有误导性。
2. 外部批量调用`GetCF`等方法时，提供一个类似`badger.DB.View(fn func(txn *Txn) error)`闭包函数包裹事务的方法是否更好