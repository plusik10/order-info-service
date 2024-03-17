package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/plusik10/cmd/order-info-service/internal/model"
	"github.com/plusik10/cmd/order-info-service/internal/repository"
	_const "github.com/plusik10/cmd/order-info-service/internal/repository/postgres/const"
	"github.com/plusik10/cmd/order-info-service/pkg/db"
)

type orderRepository struct {
	client db.Client
}

func NewOrderRepository(client db.Client) repository.OrderRepository {
	return &orderRepository{client: client}
}

// Create implements repository.Repository.
func (r *orderRepository) Create(ctx context.Context, order model.Order) error {
	op := "repository.OrderRepository.Create"
	err := r.CreateOrder(ctx, order)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	err = r.CreateDelivery(ctx, order.OrderUID, order.Delivery)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	err = r.CreatePayment(ctx, order.OrderUID, order.Payment)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	for _, item := range order.Items {
		err = r.CreateItems(ctx, order.OrderUID, item)
		if err != nil {
			return errors.New(op + ": " + err.Error())
		}
	}

	return nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, order model.Order) error {
	op := "repository.OrderRepository.CreateOrder"
	query, arg, err := squirrel.
		Insert(_const.ORDER_TABLE).
		Columns(_const.ORDER_UID,
			_const.ENTRY,
			_const.TRACK_NUMBER,
			_const.LOCALE,
			_const.INTERNAL_SIGNATURE,
			_const.CUSTOMER_ID,
			_const.DELIVERY_SERVICE,
			_const.SHARD_KEY,
			_const.SM_ID,
			_const.DATE_CREATED,
			_const.OOF_SHARD).
		Values(order.OrderUID,
			order.Entry,
			order.TrackNumber,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.Shardkey,
			order.SmID,
			order.DateCreated,
			order.OofShard).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	q := db.Query{
		Name:     "CreateOrder",
		QueryRow: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, arg...)

	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	return nil
}

func (r *orderRepository) CreatePayment(ctx context.Context, orderUid string, payment model.Payment) error {
	op := "repository.OrderRepository.CreatePayment"
	query, arg, err := squirrel.
		Insert(_const.PAYMENT_TABLE).
		Columns(_const.ORDER_UID,
			_const.TRANSACTION,
			_const.REQUEST_ID,
			_const.CURRENCY,
			_const.PROVIDER,
			_const.AMOUNT,
			_const.PAYMENT_DT,
			_const.BANK,
			_const.DELIVERY_COST,
			_const.GOODS_TOTAL,
			_const.CUSTOM_FEE).
		Values(orderUid,
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDt,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	q := db.Query{
		Name:     "CreatePayment",
		QueryRow: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, arg...)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	return nil
}

func (r *orderRepository) CreateDelivery(ctx context.Context, orderUid string, delivery model.Delivery) error {
	op := "repository.OrderRepository.CreateDelivery"
	query, arg, err := squirrel.
		Insert(_const.DELIVERY_TABLE).
		Columns(_const.ORDER_UID,
			_const.NAME,
			_const.PHONE,
			_const.ZIP,
			_const.CITY,
			_const.ADDRESS,
			_const.REGION,
			_const.EMAIL).
		Values(orderUid,
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	q := db.Query{
		Name:     "CreateDelivery",
		QueryRow: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, arg...)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	return nil
}

func (r *orderRepository) CreateItems(ctx context.Context, OrderUID string, item model.Item) error {
	op := "repository.OrderRepository.CreateItems"
	query, arg, err := squirrel.
		Insert(_const.ITEMS_TABLE).
		Columns(_const.CHRT_ID,
			_const.ORDER_UID,
			_const.TRACK_NUMBER,
			_const.PRICE,
			_const.RID,
			_const.NAME,
			_const.SALE,
			_const.SIZE,
			_const.TOTAL_PRICE,
			_const.NM_ID,
			_const.BRAND,
			_const.STATUS).
		Values(item.ChrtID,
			OrderUID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	q := db.Query{
		Name:     "CreateItem",
		QueryRow: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, arg...)
	if err != nil {
		return errors.New("error: " + err.Error() + "query " + q.QueryRow)
	}

	return nil
}
