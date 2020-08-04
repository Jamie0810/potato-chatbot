package model

import "time"

const (
	// Deposit 打款
	Deposit TransactionType = "deposit"
	// Withdrawal 收款
	Withdrawal TransactionType = "withdrawal"
)

type TransactionType string

type Order struct {
	TrackingNumber      string          `json:"tarcking_number"`
	TransactionType     TransactionType `json:"transaction_type"`
	PaymentNumber       string          `json:"payment_number"`
	EWallet             string          `json:"payment_type"`
	CreatedAt           time.Time       `json:"crated_at"`
	MerchantID          uint64          `json:"merchant_id"`
	MerchantName        string          `json:"merchant_name"`
	MerchantProjectName string          `json:"merchant_project_name"`
	MerchantOrderNumber string          `json:"merchant_order_number"`
	ChannelName         string          `json:"cannel_name"`
	ChatroomIDs         []int64         `json:"chanroom_ids"`
}

type ChatRoom struct {
	ChatID    int64  `gorm:"column:chat_id;type:bigint;" json:"chatID" form:"chatID"`               //聊天室ID
	GroupName string `gorm:"column:group_name;type:varchar(30);" json:"groupName" form:"groupName"` //群組名稱
	UserID    int64  `gorm:"column:user_id;type:bigint;" json:"userID" form:"userID"`               //使用者ID
	UserName  string `gorm:"column:user_name;type:varchar(30);" json:"userName" form:"userName"`    //使用者名稱
}
