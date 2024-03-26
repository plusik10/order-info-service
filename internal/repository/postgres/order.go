package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

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
	err := r.createOrder(ctx, order)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	err = r.createDelivery(ctx, order.OrderUID, order.Delivery)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	err = r.createPayment(ctx, order.OrderUID, order.Payment)
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	for _, item := range order.Items {
		err = r.createItems(ctx, order.OrderUID, item)
		if err != nil {
			return errors.New(op + ": " + err.Error())
		}
	}

	return nil
}

func (r *orderRepository) createOrder(ctx context.Context, order model.Order) error {
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

func (r *orderRepository) createPayment(ctx context.Context, orderUid string, payment model.Payment) error {
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

func (r *orderRepository) createDelivery(ctx context.Context, orderUid string, delivery model.Delivery) error {
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

func (r *orderRepository) createItems(ctx context.Context, OrderUID string, item model.Item) error {
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

func (r *orderRepository) GetOrderUIDs(ctx context.Context) ([]string, error) {
	op := "repository.OrderRepository.GetOrderUIDs"
	result := make([]string, 0)
	query, arg, err := squirrel.Select(_const.ORDER_UID).
		From(_const.ORDER_TABLE).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {

		return result, errors.New(op + ": " + err.Error())
	}
	q := db.Query{
		Name:     "GetAll",
		QueryRow: query,
	}

	rows, err := r.client.DB().QueryContext(ctx, q, arg...)
	if err != nil {
		return result, errors.New(op + ": " + err.Error())
	}
	defer rows.Close() // Важно закрыть rows после использования

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, errors.New(op + ": " + err.Error())
		}
		result = append(result, uid)
	}

	return result, nil
}

func (r *orderRepository) getOrder(ctx context.Context, orderUID string) (model.Order, error) {
	op := "repository.OrderRepository.getOrder"
	query, arg, err := squirrel.
		Select(_const.TRACK_NUMBER,
			_const.ENTRY,
			_const.LOCALE,
			_const.INTERNAL_SIGNATURE,
			_const.CUSTOMER_ID,
			_const.DELIVERY_SERVICE,
			_const.SHARD_KEY,
			_const.SM_ID,
			_const.DATE_CREATED,
			_const.OOF_SHARD).
		From(_const.ORDER_TABLE).
		Where(squirrel.Eq{_const.ORDER_UID: orderUID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Order{}, errors.New(op + ": " + err.Error())
	}
	q := db.Query{
		Name:     "GetOrderByUID",
		QueryRow: query,
	}

	var order model.Order
	err = r.client.DB().GetContext(ctx, &order, q, arg...)

	if err != nil {
		return model.Order{}, errors.New(op + ": " + err.Error())
	}

	return order, nil
}

func (r *orderRepository) getDelivery(ctx context.Context, orderUID string) (model.Delivery, error) {
	op := "repository.OrderRepository.getDelivery"
	query, arg, err := squirrel.
		Select(_const.NAME,
			_const.PHONE,
			_const.ZIP,
			_const.CITY,
			_const.ADDRESS,
			_const.REGION,
			_const.EMAIL,
		).
		From(_const.DELIVERY_TABLE).
		Where(squirrel.Eq{_const.ORDER_UID: orderUID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		fmt.Println("Ошибка формирования запроса")
		return model.Delivery{}, errors.New(op + ": " + err.Error())
	}

	q := db.Query{
		Name:     "GetDeliveryByUID",
		QueryRow: query,
	}

	var delivery model.Delivery
	err = r.client.DB().GetContext(ctx, &delivery, q, arg...)
	if err != nil {
		// TODO: upd
		for _, a := range arg {
			fmt.Println("arg: ", a)
		}
		return model.Delivery{}, errors.New(op + ": " + err.Error())
	}

	return delivery, nil
}

func (r *orderRepository) getPayment(ctx context.Context, orderUID string) (model.Payment, error) {
	op := "repository.OrderRepository.GetPayment"
	query, arg, err := squirrel.
		Select("p."+_const.TRANSACTION,
			_const.REQUEST_ID,
			_const.CURRENCY,
			_const.PROVIDER,
			_const.AMOUNT,
			_const.PAYMENT_DT,
			_const.BANK,
			_const.DELIVERY_COST,
			_const.GOODS_TOTAL,
			_const.CUSTOM_FEE).
		From(_const.PAYMENT_TABLE + " p").
		Where(squirrel.Eq{_const.ORDER_UID: orderUID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return model.Payment{}, errors.New(op + ": " + err.Error())
	}

	q := db.Query{
		Name:     "GetPaymentByUID",
		QueryRow: query,
	}

	var payment model.Payment
	err = r.client.DB().GetContext(ctx, &payment, q, arg...)

	if err != nil {
		fmt.Println("GetPaymentByUID ошибка", err.Error())
		fmt.Println("Запрос: ", q.QueryRow)
		for _, i := range arg {
			fmt.Println(i)
		}
		return model.Payment{}, errors.New(op + ": " + err.Error())
	}

	return payment, nil
}

func (r *orderRepository) getItems(ctx context.Context, orderUID string) ([]model.Item, error) {
	op := "repository.OrderRepository.getItems"
	query, arg, err := squirrel.
		Select(_const.CHRT_ID,
			_const.TRACK_NUMBER,
			_const.PRICE,
			_const.RID,
			_const.NAME,
			_const.SALE,
			_const.SIZE,
			_const.TOTAL_PRICE,
			_const.NM_ID,
			_const.BRAND,
			_const.STATUS,
		).
		From(_const.ITEMS_TABLE).
		Where(squirrel.Eq{_const.ORDER_UID: orderUID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.New(op + ": " + err.Error())
	}

	var items []model.Item

	q := db.Query{
		Name:     "GetAllItems",
		QueryRow: query,
	}

	err = r.client.DB().SelectContext(ctx, &items, q, arg...)
	if err != nil {
		fmt.Println("ОШИБКА ТУТ ПОЛУЧАЕТСЯ")
		return nil, errors.New(op + ": " + err.Error())
	}

	return items, nil
}

func (r *orderRepository) GetOrderByUID(ctx context.Context, orderUID string) (model.Order, error) {
	op := "repository.OrderRepository.GetOrderByUID"
	var order model.Order

	order, err := r.getOrder(ctx, orderUID)
	order.OrderUID = orderUID

	if err != nil {
		return order, errors.New(op + ": " + err.Error())
	}
	delivery, err := r.getDelivery(ctx, orderUID)
	if err != nil {
		return order, errors.New(op + ": " + err.Error())
	}

	payment, err := r.getPayment(ctx, orderUID)
	if err != nil {
		return order, errors.New(op + ": " + err.Error())
	}

	items, err := r.getItems(ctx, orderUID)
	if err != nil {
		return order, errors.New(op + ": " + err.Error())
	}

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	return order, nil
}

func (r *orderRepository) GetOrders(ctx context.Context) ([]*model.Order, error) {
	//op := "repository.OrderRepository.GetOrders"

	query, arg, err := squirrel.
		Select(_const.ORDER_UID,
			_const.TRACK_NUMBER,
			_const.ENTRY,
			_const.LOCALE,
			_const.INTERNAL_SIGNATURE,
			_const.CUSTOMER_ID,
			_const.DELIVERY_SERVICE,
			_const.SHARD_KEY,
			_const.SM_ID,
			_const.DATE_CREATED,
			_const.OOF_SHARD).
		From(_const.ORDER_TABLE).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		fmt.Println("error getting orders: ", err)
	}

	var orders []*model.Order

	q := db.Query{
		Name:     "GetAll",
		QueryRow: query,
	}
	err = r.client.DB().SelectContext(ctx, &orders, q, arg...)
	if err != nil {
		log.Println("error geting orders")
		return nil, err
	}

	if len(orders) == 0 {
		return nil, nil
	}

	for i, order := range orders {
		delivey, err := r.getDelivery(ctx, order.OrderUID)
		if err != nil {
			log.Println("error getting delivery: ", err)
		}
		orders[i].Delivery = delivey

		payment, err := r.getPayment(ctx, order.OrderUID)
		if err != nil {
			fmt.Println("error getting payment: ", err)
		}
		orders[i].Payment = payment

		items, err := r.getItems(ctx, order.OrderUID)
		if err != nil {
			fmt.Println("error getting items")
		}
		orders[i].Items = items
	}

	return orders, nil
}
