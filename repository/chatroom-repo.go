package repository

import (
	"context"
	"fmt"
	"mods/entity"

	"gorm.io/gorm"
)

type chatroomConnection struct {
	connection *gorm.DB
}

type ChatroomRepository interface {
	// functional
	AddChatroom(ctx context.Context, chatroom entity.ChatRoom) (entity.ChatRoom, error)
	RemoveChatroom(ctx context.Context, id uint64) error
	GetChatroomUser(ctx context.Context, id string) ([]entity.ChatRoom, error)
	GetChatroomDoctor(ctx context.Context, id string) ([]entity.ChatRoom, error)
	IsDuplicateChatRoom(ctx context.Context, U_id string , D_id string) (bool, error)
}

func NewChatroomRepository(db *gorm.DB) ChatroomRepository {
	return &chatroomConnection{
		connection: db,
	}
}

func (db *chatroomConnection) AddChatroom(ctx context.Context, chatroom entity.ChatRoom) (entity.ChatRoom, error) {
	if err := db.connection.Create(&chatroom).Error; err != nil {
		return entity.ChatRoom{}, err
	}

	return chatroom, nil
}

func (db *chatroomConnection) RemoveChatroom(ctx context.Context, id uint64) error {
	var chatroom entity.ChatRoom
	
	tx := db.connection.Where("id = ?", id).Delete(&chatroom)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := fmt.Errorf("no record found with id: %v", id)
		return err
	}

	return nil
}

func (db *chatroomConnection) IsDuplicateChatRoom(ctx context.Context, U_id string , D_id string) (bool, error) {
	var chatroom entity.ChatRoom
	tx := db.connection.Where("uid = ?", U_id).Where("uid_doctor = ?", D_id).First(&chatroom)

	if tx.Error != nil {
		return false, nil
	}

	return true, nil
}

func (db *chatroomConnection) GetChatroomUser(ctx context.Context, id string) ([]entity.ChatRoom, error) {
	var listChatroom []entity.ChatRoom
	// var doctor = "Doctor"

	tx := db.connection.Where("uid = ?", id).Find(&listChatroom)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return listChatroom, nil
}

func (db *chatroomConnection) GetChatroomDoctor(ctx context.Context, id string) ([]entity.ChatRoom, error) {
	var listChatroom []entity.ChatRoom
	// var doctor = "Doctor"

	tx := db.connection.Where("uid_doctor = ?", id).Find(&listChatroom)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return listChatroom, nil
}