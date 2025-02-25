/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azurequeuestoragesource

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"

	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	"github.com/triggermesh/triggermesh/pkg/apis/sources/v1alpha1"
)

// envConfig is a set parameters sourced from the environment for the source's
// adapter.
type envConfig struct {
	pkgadapter.EnvConfig

	AccountName string `envconfig:"AZURE_ACCOUNT_NAME"`
	AccountKey  string `envconfig:"AZURE_ACCOUNT_KEY"`
	QueueName   string `envconfig:"AZURE_QUEUE_NAME"`
}

// NewEnvConfig returns an accessor for the source's adapter envConfig.
func NewEnvConfig() pkgadapter.EnvConfigAccessor {
	return &envConfig{}
}

const (
	concurrentMsgProcessing = 16
)

// adapter implements the source's adapter.
type adapter struct {
	messagesURL azqueue.MessagesURL
	ceClient    cloudevents.Client
	eventsource string
	logger      *zap.SugaredLogger
}

// NewAdapter returns a constructor for the source's adapter.
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envConfig)
	logger := logging.FromContext(ctx)

	credential, err := azqueue.NewSharedKeyCredential(env.AccountName, env.AccountKey)
	if err != nil {
		logger.Fatal("azqueue.NewSharedKeyCredential failed: ", err)
	}

	p := azqueue.NewPipeline(credential, azqueue.PipelineOptions{})

	urlRef, err := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net", env.AccountName))

	if err != nil {
		logger.Fatal("url.Parse failed: ", err)
	}

	serviceURL := azqueue.NewServiceURL(*urlRef, p)

	// Create a URL that references a queue in your Azure Storage account.
	// This returns a QueueURL object that wraps the queue's URL and a request pipeline (inherited from serviceURL)
	queueURL := serviceURL.NewQueueURL(env.QueueName) // Queue names require lowercase
	// Create a URL allowing you to manipulate a queue's messages.
	// This returns a MessagesURL object that wraps the queue's messages URL and a request pipeline (inherited from queueURL)
	messagesURL := queueURL.NewMessagesURL()

	return &adapter{
		messagesURL: messagesURL,
		ceClient:    ceClient,
		eventsource: queueURL.String(),
		logger:      logging.FromContext(ctx),
	}
}

var _ pkgadapter.Adapter = (*adapter)(nil)

// Start implements adapter.Adapter.
func (h *adapter) Start(ctx context.Context) error {
	msgCh := make(chan *azqueue.DequeuedMessage, concurrentMsgProcessing)
	h.processQueueEvents(ctx, msgCh)

	<-ctx.Done()

	return ctx.Err()

}

func (h *adapter) processQueueEvents(ctx context.Context, msgCh chan *azqueue.DequeuedMessage) {
	h.logger.Info("Starting to process queue events")
	for {
		ctx := context.TODO()
		// Create goroutines that can process messages in parallel
		dmr, err := h.messagesURL.Dequeue(ctx, int32(concurrentMsgProcessing), 24*time.Hour)
		if err != nil {
			h.logger.Error(err)
			continue
		}

		for n := 0; n < int(dmr.NumMessages()); n++ {
			message := dmr.Message(int32(n))
			msgCh <- message
			go func(msgCh <-chan *azqueue.DequeuedMessage) {

				msg := <-msgCh // Get a message from the channel
				h.logger.Infof("New msg [%v] in the channel!", msg.Text)
				// Create a URL allowing you to manipulate this message.
				// This returns a MessageIDURL object that wraps the this message's URL and a request pipeline (inherited from messagesURL)
				msgIDURL := h.messagesURL.NewMessageIDURL(msg.ID)
				popReceipt := msg.PopReceipt // This message's most-recent pop receipt

				update, err := msgIDURL.Update(ctx, popReceipt, time.Second*20, msg.Text)
				if err != nil {
					h.logger.Error(err)
					return
				}
				popReceipt = update.PopReceipt // Performing any operation on a message ID always requires the most recent pop receipt

				err = h.sendCloudEvent(msg)
				if err != nil {
					h.logger.Error(err)
					return
				}

				// After processing the message, delete it from the queue so it won't be dequeued ever again:
				_, err = msgIDURL.Delete(ctx, popReceipt)
				if err != nil {
					h.logger.Error(err)
					return
				}
			}(msgCh)
		}
	}
}

func (h *adapter) sendCloudEvent(m *azqueue.DequeuedMessage) error {
	event := cloudevents.NewEvent(cloudevents.VersionV1)
	event.SetType(v1alpha1.AzureQueueStorageEventType)
	event.SetSource(h.eventsource)
	if err := event.SetData(cloudevents.ApplicationJSON, m); err != nil {
		return fmt.Errorf("failed to set event data: %w", err)
	}

	h.logger.Debug("Sending CloudEvent: ", event)
	if result := h.ceClient.Send(context.Background(), event); !cloudevents.IsACK(result) {
		return fmt.Errorf("failed to send CloudEvent: %w", result)
	}
	return nil
}
