// +build wasm

package ipc

import (
	"encoding/binary"
	"syscall/js"
)

type msgChunk struct {
	Channel    string
	MsgId      int32
	ThisChunk  int32
	ChunkCount int32
	Content    []byte
}

func (m msgChunk) serialize() []byte {
	chanBytes := []byte(m.Channel)
	result := make([]byte, 2+len(chanBytes)+len(m.Content)+3*4)
	binary.BigEndian.PutUint16(result[0:2], uint16(len(chanBytes)))
	binary.BigEndian.PutUint32(result[2:6], uint32(m.MsgId))
	binary.BigEndian.PutUint32(result[6:10], uint32(m.ThisChunk))
	binary.BigEndian.PutUint32(result[10:14], uint32(m.ChunkCount))
	copy(result[14:14+len(chanBytes)], chanBytes)
	copy(result[14+len(chanBytes):], m.Content)
	return result
}

func deserializeMsgChunk(data []byte) msgChunk {
	result := msgChunk{}
	chanLen := int(binary.BigEndian.Uint16(data[0:2]))
	result.MsgId = int32(binary.BigEndian.Uint32(data[2:6]))
	result.ThisChunk = int32(binary.BigEndian.Uint32(data[6:10]))
	result.ChunkCount = int32(binary.BigEndian.Uint32(data[10:14]))
	result.Channel = string(data[14 : 14+chanLen])
	result.Content = data[14+chanLen:]
	return result
}

const channelChunk = "channelChunk"
const chunkSize = 4096

var msgIdGen = int32(0)

// SendChunked sends a message to a specific sender in a chunked format.  This allows the sending of enormous
// messages without overruning potentially small message size limits.  Receiving a chunked message requires wrapping
// a chunked endpoint around the receiving end of the endpoint, via GetChunkedEndpoint(...)
func SendChunked(sender Sender, channel string, msg []byte) {
	msgId := msgIdGen
	msgIdGen++
	chunkCount := (len(msg) + chunkSize - 1) / chunkSize
	if chunkCount == 0 {
		chunkCount = 1
	}
	var msgs [][]byte
	for i := 0; i < chunkCount; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(msg) {
			end = len(msg)
		}
		chunk := msgChunk{
			Channel:    channel,
			MsgId:      msgId,
			ChunkCount: int32(chunkCount),
			ThisChunk:  int32(i),
			Content:    msg[start:end],
		}
		msgs = append(msgs, chunk.serialize())
	}
	for _, msg := range msgs {
		sender.Send(channelChunk, msg)
	}
}

// GetChunkedEndpoint creates an wrapper around the specified endpoint facilitating the sending and receipt of
// chunked messages
func GetChunkedEndpoint(targetEndpoint Endpoint) Endpoint {
	if _, ok := endpointsWithChunkedListeners[targetEndpoint]; !ok {
		endpointsWithChunkedListeners[targetEndpoint] = newChunkedEndpoint(targetEndpoint)
	}
	return endpointsWithChunkedListeners[targetEndpoint]
}

type chunkedEndpoint struct {
	endpoint Endpoint
	partials map[int32]map[int32][]byte
	outChans map[string]map[int]chan<- Message
}

func newChunkedEndpoint(endpoint Endpoint) *chunkedEndpoint {
	cl := &chunkedEndpoint{
		endpoint: endpoint,
		partials: map[int32]map[int32][]byte{},
		outChans: map[string]map[int]chan<- Message{},
	}
	inChunks, _ := endpoint.Listen(channelChunk)
	go func() {
		for {
			chunkRawMsg := <-inChunks
			chunkMsg := deserializeMsgChunk(chunkRawMsg.Content())
			msgPartials := cl.getPartial(chunkMsg.MsgId)
			msgPartials[chunkMsg.ThisChunk] = chunkMsg.Content
			if len(msgPartials) != int(chunkMsg.ChunkCount) {
				continue
			}
			delete(cl.partials, chunkMsg.MsgId)
			var fullMsgBytes []byte
			for i := int32(0); i < chunkMsg.ChunkCount; i++ {
				fullMsgBytes = append(fullMsgBytes, msgPartials[i]...)
			}
			outMsg := &chunkedMessage{
				source:  cl,
				content: fullMsgBytes,
			}
			if listeners, ok := cl.outChans[chunkMsg.Channel]; ok {
				for _, outChan := range listeners {
					func() {
						defer func() { recover() }()
						outChan <- outMsg
					}()
				}
			}
		}
	}()
	return cl
}

func (ce *chunkedEndpoint) Send(channelName string, data []byte) {
	SendChunked(ce.endpoint, channelName, data)
}

var listenerHandleGen = 0

func (ce *chunkedEndpoint) Listen(channelName string) (recv <-chan Message, done func()) {
	msgChan := make(chan Message)
	if _, ok := ce.outChans[channelName]; !ok {
		ce.outChans[channelName] = map[int]chan<- Message{}
	}
	listenerId := listenerHandleGen
	listenerHandleGen++
	ce.outChans[channelName][listenerId] = msgChan
	doneFunc := func() {
		delete(ce.outChans[channelName], listenerId)
	}
	return msgChan, doneFunc
}

func (ce *chunkedEndpoint) getPartial(msgId int32) map[int32][]byte {
	update, ok := ce.partials[msgId]
	if !ok {
		update = map[int32][]byte{}
		ce.partials[msgId] = update
	}
	return update
}

type chunkedMessage struct {
	source  *chunkedEndpoint
	content []byte
}

func (c chunkedMessage) Content() []byte { return c.content }
func (c chunkedMessage) Sender() string  { return "chunked" }
func (c chunkedMessage) Reply(channelName string, data []byte) {
}

func (c chunkedMessage) JSValue() js.Value { return js.Undefined() }

var endpointsWithChunkedListeners = map[Endpoint]*chunkedEndpoint{}
