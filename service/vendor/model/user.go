package model

import (
	"entity"

	"util"
	log "util/logger"
)

func init() {
	u := &entity.UserInfoSerializable{}
	if !agendaDB.HasTable(u) {
		err := agendaDB.CreateTable(u).Error
		util.PanicIf(err)
		log.Infof("\n ...... CreateTable %T. \n", u)
	}
}

// UserInfoAtomicService .
type UserInfoAtomicService struct{}

// UserInfoService .
var UserInfoService = UserInfoAtomicService{}

// func loadAllUserFromDB(db *DB) {}

// Create .
func (*UserInfoAtomicService) Create(u *entity.UserInfoSerializable) error {
	tx := agendaDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Save .
func (*UserInfoAtomicService) Save(u *entity.UserInfoSerializable) error {
	tx := agendaDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Save(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete .
func (*UserInfoAtomicService) Delete(u *entity.UserInfoSerializable) error {
	tx := agendaDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Delete(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// FindAll .
func (*UserInfoAtomicService) FindAll() []entity.UserInfoSerializable {
	var rows []entity.UserInfoSerializable
	agendaDB.Find(&rows) // CHECK: should check .Error ?
	return rows
}

// FindByUsername .
func (*UserInfoAtomicService) FindByUsername(name entity.Username) *entity.UserInfoSerializable {
	var uInfo entity.UserInfoSerializable

	// agendaDB.First(&uInfo, entity.UserInfoSerializable{Name: name}) // TODEL: sad to anonymous member ...
	u := entity.UserInfoSerializable{}
	u.Name = name
	agendaDB.First(&uInfo, u)

	return &uInfo
}
