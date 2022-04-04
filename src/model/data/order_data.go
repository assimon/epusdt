package data

import (
	"context"
	"fmt"
	"github.com/assimon/luuu/model/dao"
	"github.com/assimon/luuu/model/mdb"
	"github.com/assimon/luuu/model/request"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	CacheWalletAddressToAvailableAmountKey = "WalletAddressToAvailableAmountKey:" // 钱包 - 可用金额 - 过期时间
	CacheWalletAddressToOrdersKey          = "WalletAddressToOrdersKey:"          // 钱包 - 待支付金额 - 订单号
)

// GetOrderInfoByOrderId 通过客户订单号查询订单
func GetOrderInfoByOrderId(orderId string) (*mdb.Orders, error) {
	order := new(mdb.Orders)
	err := dao.Mdb.Model(order).Limit(1).Find(order, "order_id = ?", orderId).Error
	return order, err
}

// GetOrderInfoByTradeId 通过交易号查询订单
func GetOrderInfoByTradeId(tradeId string) (*mdb.Orders, error) {
	order := new(mdb.Orders)
	err := dao.Mdb.Model(order).Limit(1).Find(order, "trade_id = ?", tradeId).Error
	return order, err
}

// CreateOrderWithTransaction 事务创建订单
func CreateOrderWithTransaction(tx *gorm.DB, order *mdb.Orders) error {
	err := tx.Model(order).Create(order).Error
	return err
}

// GetOrderByBlockIdWithTransaction 通过区块获取订单
func GetOrderByBlockIdWithTransaction(tx *gorm.DB, blockId string) (*mdb.Orders, error) {
	order := new(mdb.Orders)
	err := tx.Model(order).Limit(1).Find(order, "block_transaction_id = ?", blockId).Error
	return order, err
}

// OrderSuccessWithTransaction 事务支付成功
func OrderSuccessWithTransaction(tx *gorm.DB, req *request.OrderProcessingRequest) error {
	err := tx.Model(&mdb.Orders{}).Where("trade_id = ?", req.TradeId).Updates(map[string]interface{}{
		"block_transaction_id": req.BlockTransactionId,
		"status":               mdb.StatusPaySuccess,
		"callback_confirm":     mdb.CallBackConfirmNo,
	}).Error
	return err
}

// GetPendingCallbackOrders 查询出等待回调的订单
func GetPendingCallbackOrders() ([]mdb.Orders, error) {
	var orders []mdb.Orders
	err := dao.Mdb.Model(orders).
		Where("callback_num < ?", 5).
		Where("callback_confirm = ?", mdb.CallBackConfirmNo).
		Where("status = ?", mdb.StatusPaySuccess).
		Find(&orders).Error
	return orders, err
}

// SaveCallBackOrdersResp 保存订单回调结果
func SaveCallBackOrdersResp(order *mdb.Orders) error {
	err := dao.Mdb.Model(order).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"callback_num":     gorm.Expr("callback_num + ?", 1),
		"callback_confirm": order.CallBackConfirm,
	}).Error
	return err
}

// UpdateOrderIsExpirationById 通过id设置订单过期
func UpdateOrderIsExpirationById(id uint64) error {
	err := dao.Mdb.Model(mdb.Orders{}).Where("id = ?", id).Update("status", mdb.StatusExpired).Error
	return err
}

// LockPayCache 锁定支付池
func LockPayCache(token, tradeId, amount string, expirationTime int64) error {
	// 载入Redis
	ctx := context.Background()
	pipe := dao.Rdb.TxPipeline()
	lockAvailableAmountKey := fmt.Sprintf("%s%s", CacheWalletAddressToAvailableAmountKey, token)
	// 占用金额
	pipe.HSet(ctx, lockAvailableAmountKey, amount, expirationTime)
	// 标记金额对应订单号
	lockWalletOrdersKey := fmt.Sprintf("%s%s", CacheWalletAddressToOrdersKey, token)
	pipe.HSet(ctx, lockWalletOrdersKey, amount, tradeId)
	_, err := pipe.Exec(ctx)
	return err
}

// ClearPayCache 清理支付池
func ClearPayCache(token, amount string) error {
	ctx := context.Background()
	lockAvailableAmountKey := fmt.Sprintf("%s%s", CacheWalletAddressToAvailableAmountKey, token)
	lockWalletOrdersKey := fmt.Sprintf("%s%s", CacheWalletAddressToOrdersKey, token)
	pipe := dao.Rdb.TxPipeline()
	pipe.HDel(ctx, lockAvailableAmountKey, amount)
	pipe.HDel(ctx, lockWalletOrdersKey, amount)
	_, err := pipe.Exec(ctx)
	return err
}

func GetTradeIdByAmount(ctx context.Context, token, amount string) (string, error) {
	cacheKey := fmt.Sprintf("%s%s", CacheWalletAddressToOrdersKey, token)
	result, err := dao.Rdb.HGet(ctx, cacheKey, amount).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

func GetExpirationTimeByAmount(ctx context.Context, token, amount string) (string, error) {
	cacheKey := fmt.Sprintf("%s%s", CacheWalletAddressToAvailableAmountKey, token)
	result, err := dao.Rdb.HGet(ctx, cacheKey, amount).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}
