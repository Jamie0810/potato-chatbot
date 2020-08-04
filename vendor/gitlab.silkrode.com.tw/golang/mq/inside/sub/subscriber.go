package sub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"log"
)

const (
	syncMode                      = "sync"
	asyncMode                     = "async"
	defaultMaxOutstandingMessages = 20
)

type Subscriber struct {
	ctx                context.Context
	subscription       *pubsub.Subscription
	pullMode           string
	defaultErrorHandle func(err error, msg *pubsub.Message)
	errorHook          func(err error, msgData string, msgID string)
}

func NewSubscriber(ctx context.Context, subscription *pubsub.Subscription) *Subscriber {
	subscription.ReceiveSettings.MaxOutstandingMessages = defaultMaxOutstandingMessages
	return &Subscriber{
		ctx:          ctx,
		subscription: subscription,
		defaultErrorHandle: func(err error, msg *pubsub.Message) {
			errorLog := ""
			if msg != nil {
				errorLog = fmt.Sprintf("mq err: %s,msgData: %s,msgID: %s", err.Error(), msg.Data, msg.ID)
			} else {
				errorLog = fmt.Sprintf("mq err: %s", err.Error())
			}
			log.Print(errorLog)
		},
	}
}

func (sb *Subscriber) Options(opts ...SubscriberOption) *Subscriber {
	for _, opt := range opts {
		opt(sb)
	}
	return sb
}

func (sb *Subscriber) Subscribe(process func(context.Context, []byte, string) error) {
	var err error
	if sb.pullMode == syncMode {
		sb.subscription.ReceiveSettings.Synchronous = true
		err = sb.pullMessageSync(sb.ctx, process)
	}
	if sb.pullMode == asyncMode {
		sb.subscription.ReceiveSettings.Synchronous = false
		err = sb.pullMessageAsync(sb.ctx, process)
	}
	if err != nil {
		if sb.errorHook != nil {
			sb.errorHook(err, "", "")
		} else {
			sb.defaultErrorHandle(err, nil)
		}
	}
}

func (sb *Subscriber) pullMessageSync(ctx context.Context, process func(context.Context, []byte, string) error) error {
	cm := make(chan *pubsub.Message)
	go func() {
		for {
			select {
			case msg := <-cm:
				err := process(sb.ctx, msg.Data, msg.ID)
				if err != nil {
					if sb.errorHook != nil {
						sb.errorHook(err, string(msg.Data), msg.ID)
					} else {
						sb.defaultErrorHandle(err, msg)
					}
				}
				msg.Ack()
			case <-ctx.Done():
				return
			}
		}
	}()
	err := sb.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	return err
}

func (sb *Subscriber) pullMessageAsync(ctx context.Context, process func(context.Context, []byte, string) error) error {
	err := sb.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		err := process(sb.ctx, msg.Data, msg.ID)
		if err != nil {
			if sb.errorHook != nil {
				sb.errorHook(err, string(msg.Data), msg.ID)
			} else {
				sb.defaultErrorHandle(err, msg)
			}
		}
		msg.Ack()
	})
	return err
}
