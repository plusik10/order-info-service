package orderRepositoryConst

const (
	// order_doc
	ORDER_TABLE        = "order_doc"
	ORDER_UID          = "order_uid"
	TRACK_NUMBER       = "track_number"
	ENTRY              = "entry"
	LOCALE             = "locale"
	INTERNAL_SIGNATURE = "internal_signature"
	CUSTOMER_ID        = "customer_id"
	DELIVERY_SERVICE   = "delivery_service"
	SHARD_KEY          = "shardkey"
	SM_ID              = "sm_id"
	DATE_CREATED       = "date_created"
	OOF_SHARD          = "oof_shard"

	// payment
	PAYMENT_TABLE = "payment"
	TRANSACTION   = "transaction"
	REQUEST_ID    = "request_id"
	CURRENCY      = "currency"
	PROVIDER      = "provider"
	AMOUNT        = "amount"
	PAYMENT_DT    = "payment_dt"
	BANK          = "bank"
	DELIVERY_COST = "delivery_cost"
	GOODS_TOTAL   = "goods_total"
	CUSTOM_FEE    = "custom_fee"

	// delivery
	DELIVERY_TABLE = "delivery"
	NAME           = "name"
	PHONE          = "phone"
	ZIP            = "zip"
	CITY           = "city"
	ADDRESS        = "address"
	REGION         = "region"
	EMAIL          = "email"

	// items
	ITEMS_TABLE = "item"
	CHRT_ID     = "chrt_id"
	PRICE       = "price"
	RID         = "rid"
	SALE        = "sale"
	SIZE        = "size"
	TOTAL_PRICE = "total_price"
	NM_ID       = "nm_id"
	BRAND       = "brand"
	STATUS      = "status"
)
