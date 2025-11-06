
package models

import "time"

type SlotStatus string
type SwapStatus string

const (
	SlotBusy       SlotStatus = "BUSY"
	SlotSwappable  SlotStatus = "SWAPPABLE"
	SlotSwapPending SlotStatus = "SWAP_PENDING"

	SwapPending  SwapStatus = "PENDING"
	SwapAccepted SwapStatus = "ACCEPTED"
	SwapRejected SwapStatus = "REJECTED"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string `gorm:"size:200;not null"`
	Email        string `gorm:"size:200;uniqueIndex;not null"`
	Password string `gorm:"size:300;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Event struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Title     string     `gorm:"size:300;not null"`
	StartTime time.Time  `gorm:"not null"`
	EndTime   time.Time  `gorm:"not null"`
	Status    SlotStatus `gorm:"type:VARCHAR(20);not null;default:'BUSY'"`
	UserID   uint		 `gorm:"userId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SwapRequest struct {
	
	ID           uint        `gorm:"primaryKey" json:"id"`
	MySlotID     uint        `json:"mySlotId"`
	TheirSlotID  uint        `json:"theirSlotId"`
	RequesterID  uint        `json:"requesterId"`
	ReceiverID   uint        `json:"receiverId"`
	Status       SwapStatus  `gorm:"type:VARCHAR(20);not null;default:'PENDING'"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}
