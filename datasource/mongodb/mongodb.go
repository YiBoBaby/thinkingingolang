package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// 建立连接
func BuildMongoDBConnect(uri string) (*mongo.Client, error) {
	// 连接参数
	clientOpts := options.Client()
	clientOpts.ApplyURI(uri)
	clientOpts.SetMaxPoolSize(5)
	clientOpts.SetMinPoolSize(1)
	// 连接池事件监听
	clientOpts.SetPoolMonitor(&event.PoolMonitor{
		Event: func(event *event.PoolEvent) {
			printBlue(fmt.Sprintf("PoolMonitor event:%v \n", event))
		},
	})
	// 命令行监视器
	clientOpts.SetMonitor(&event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			printGreen(fmt.Sprintf("CommandMonitor Started: ctx:%v, event:%v \n", ctx, event))
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			printGreen(fmt.Sprintf("CommandMonitor Succeeded: ctx:%v, event:%v \n", ctx, event))
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			printGreen(fmt.Sprintf("CommandMonitor Failed: ctx:%v, event:%v \n", ctx, event))
		},
	})

	// 创建客户端
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, err
	}

	// 使用上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 后台建立连接, 需要使用 ping 校验连接是否正确
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Ping
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
