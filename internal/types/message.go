/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package types

func NewMessage(senderID, receiverID, data []byte) *Message {
	return &Message{
		SenderId:   senderID,
		ReceiverId: receiverID,
		Data:       data,
	}
}
