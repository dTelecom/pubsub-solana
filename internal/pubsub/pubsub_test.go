package pubsub_test

import (
	"context"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/golang/mock/gomock"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/dTelecom/pubsub-solana/internal/contract_client"
	"github.com/dTelecom/pubsub-solana/internal/pubsub"
	"github.com/dTelecom/pubsub-solana/internal/pubsub/mocks"
)

type dataEncoderStub struct{}

func (*dataEncoderStub) Encode(data []byte) ([]byte, error) {
	var buffer []byte
	buffer = append(buffer, []byte("encoded")...)
	buffer = append(buffer, data...)
	return buffer, nil
}

func (*dataEncoderStub) Decode(data []byte) ([]byte, error) {
	return data[len("encoded"):], nil
}

var dataEncoder dataEncoderStub

type msgType struct {
	ID    string `borsh:"id"`
	Topic string `borsh:"topic"`
	Value []byte `borsh:"value"`
}

func mustSerialize(id string, topic string, value []byte) (result []byte) {
	content, err := borsh.Serialize(msgType{
		ID:    id,
		Topic: topic,
		Value: value,
	})
	if err != nil {
		panic(err)
	}

	result, _ = dataEncoder.Encode(content)
	return
}

func Test_Happy(t *testing.T) {
	currentNode, _ := solana.NewRandomPrivateKey()
	anotherNode, _ := solana.NewRandomPrivateKey()

	ctrl := gomock.NewController(t)
	solanaClient := mocks.NewMockContractClient(ctrl)
	magicblockClient := mocks.NewMockContractClient(ctrl)
	messageIdGenerator := mocks.NewMockMessageIdGenerator(ctrl)
	ps := pubsub.New(solanaClient, magicblockClient, messageIdGenerator, &dataEncoder)

	var incomingMessageHandler func(context.Context, contract_client.MessageData)
	var outgoingMessageHandler func(context.Context, contract_client.MessageData)

	t.Run("start", func(t *testing.T) {
		solanaClient.EXPECT().
			IsSigner(currentNode.PublicKey()).
			Times(1).
			Return(true)

		solanaClient.EXPECT().
			IsSigner(anotherNode.PublicKey()).
			Times(1).
			Return(false)

		magicblockClient.EXPECT().
			IncomingMessageSubscribe(gomock.Any(), anotherNode.PublicKey(), gomock.Any()).
			Times(1).
			DoAndReturn(func(_ context.Context, _ solana.PublicKey, handler func(context.Context, contract_client.MessageData)) error {
				incomingMessageHandler = handler
				return nil
			})

		magicblockClient.EXPECT().
			OutgoingMessageSubscribe(gomock.Any(), anotherNode.PublicKey(), gomock.Any()).
			Times(1).
			DoAndReturn(func(_ context.Context, _ solana.PublicKey, handler func(context.Context, contract_client.MessageData)) error {
				outgoingMessageHandler = handler
				return nil
			})

		err := ps.Start(context.Background(), []solana.PublicKey{currentNode.PublicKey(), anotherNode.PublicKey()})
		require.NoError(t, err, "start error")
	})

	t.Run("subscribe", func(t *testing.T) {
		ch := make(chan pubsub.Event, 1)
		ps.Subscribe("test-topic", func(ctx context.Context, event pubsub.Event) {
			ch <- event
		})

		magicblockClient.EXPECT().
			MarkAsRead(gomock.Any(), anotherNode.PublicKey(), int64(1)).
			Times(1).
			Return(solana.Signature{}, nil)

		incomingMessageHandler(context.Background(), contract_client.MessageData{
			Read:      false,
			TimeStamp: 1,
			Content:   mustSerialize("test-id", "test-topic", []byte("test-value")),
		})

		select {
		case event := <-ch:
			require.Equal(t, "test-id", event.ID)
			require.Equal(t, anotherNode.PublicKey().String(), event.FromPeerId)
			require.Equal(t, []byte("test-value"), event.Message)
		case <-time.After(time.Millisecond):
			require.FailNow(t, "timeout")
		}
	})

	t.Run("publish", func(t *testing.T) {
		messageIdGenerator.EXPECT().Generate().Times(1).Return("test-id1")

		magicblockClient.EXPECT().
			SendMessage(gomock.Any(), anotherNode.PublicKey(), mustSerialize("test-id1", "test-topic", []byte("test-value1"))).
			Times(1).
			Return(solana.Signature{}, nil)

		_, err := ps.Publish(context.Background(), "test-topic", []byte("test-value1"))
		require.NoError(t, err, "publish error")

		messageIdGenerator.EXPECT().Generate().Times(1).Return("test-id2")

		_, err = ps.Publish(context.Background(), "test-topic", []byte("test-value2"))
		require.NoError(t, err, "publish error")

		magicblockClient.EXPECT().
			SendMessage(gomock.Any(), anotherNode.PublicKey(), mustSerialize("test-id2", "test-topic", []byte("test-value2"))).
			Times(1).
			Return(solana.Signature{}, nil)

		outgoingMessageHandler(context.Background(), contract_client.MessageData{
			Read: true,
		})

		outgoingMessageHandler(context.Background(), contract_client.MessageData{
			Read: true,
		})

		messageIdGenerator.EXPECT().Generate().Times(1).Return("test-id3")

		magicblockClient.EXPECT().
			SendMessage(gomock.Any(), anotherNode.PublicKey(), mustSerialize("test-id3", "test-topic", []byte("test-value3"))).
			Times(1).
			Return(solana.Signature{}, nil)

		_, err = ps.Publish(context.Background(), "test-topic", []byte("test-value3"))
		require.NoError(t, err, "publish error")
	})
}
