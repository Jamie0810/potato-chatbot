# MQ



## Example

**mq 初始化**

```
mqInstance, err := mq.Init(ctx, topicID, credentialsFile, projectID, mq.InitPub())
if err != nil {
    log.Printf("failed to create client: %v", err)
}
```

**pub sample**

```
mqInstance.Publisher().
		Options(
			pub.SetErrorHook(func(err error, requestID string) {
				log.Printf("failed to publish requestID: %v ; message: %v", requestID, err)
			}),
		).
		PublishJSON("hellow word", time.Now().Format(time.RFC3339))
```

**pub options**
- SetErrorHook(func(err error,requestID string)) // 設定 error hoot，requestID 自定義，trace 方便。

**sub sample**

```
mqInstance.Subscriber().
		Options(
			sub.SyncMode(),
			sub.SetErrorHook(func(err error, msgData string, msgID string) {
				log.Printf("failed to publish msgID: %v ; message: %v", msgID, err)
			}),
		).
		Subscribe(func(ctx context.Context, payload []byte, s string) error {
			log.Println("Got Message ID:", string(s))
			log.Println("Got Payload:", string(payload))

			return nil
		})
```
**sub options**
- SyncMode() // 同步。
- AsyncMode() // 非同步。
- SetMaxOutstandingMessages(i int) // 設定從topic 一次拉取的訊息數量。
- SetNumGoroutines(i int) // 設定async routing 數量。
- SetErrorHook(func(err error)) // 設定錯誤訊息hook。



## Contributors

- [frankie](https://gitlab.silkrode.com.tw/frankie)