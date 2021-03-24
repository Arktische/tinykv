# tinykv课程思路
## project 1
1. `kv/storage/standalone_storage/standalone_storage.go`
    * 设计`StanAloneStorage`的结构：使用`package engine_util`提供的`Engines`类型即可。
    * `Reader`和`Write`方法实现：
        * Reader：增加一个StandAloneStorageReader结构，实现Reader接口，保存*badger.Txn。在Reader内开启只读事务，于Close中调用Discard
        * Write：使用engine_util中WriteBatch结构提供的方法
2. `kv/server/server.go`
    * 完成Raw(Put/Delete/Scan/Get)方法
<!-- ### 疑问
1. 参照`kv/storage/raft_storage/region_reader.go`中写法，`badger.Txn.Discard()`方法只在`Reader`的`Close()`方法中被调用，而测试代码`server_test.go`中并没有调用`Close()`去提交事务，测试代码这样写虽然没问题，但是有误导性。
2. 外部批量调用`GetCF`等方法时，提供一个类似`badger.DB.View(fn func(txn *Txn) error)`闭包函数包裹事务的方法是否更好
3. `region_reader.go:27`为何要单独对`badger.ErrKeyNotFound`做判断，且该情况返回错误设为`nil`?
4. WriteBatch中保存的SafePoint意义是什么呢，engine层的badger已经支持了事务回滚?
5. RawScanResponse中的Error字段缺乏明确说明，不知道应该如何返回错误 -->

## project 2
1. 
<!-- ### 疑问
1. RaftLog中committed和stabled区别是什么？ -->