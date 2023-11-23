/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package types

func NewMessage() *Message {
	return &Message{
		SenderId:   nil,
		ReceiverId: nil,
		Data:       nil,
	}
}
