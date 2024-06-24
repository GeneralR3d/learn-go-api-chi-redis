package order

import (

	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/generalr3d/learn-go-api-chi-redis/model"
	
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order: %d",id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error{	// can fail incase of failed communication or duplicate ID, so need to return error
	/*
	*	Note for receiver functions in Golang, must take in pointer to object else the function/method will operate on a COPY of that object
	*	However Golang automatically deferences the pointer in the function. Tricky!
	*
	*/
	data, err := json.Marshal(order)	//	encoding order struct into byte array
	if err != nil{
		return fmt.Error("Failed to encode order: %w", err)
	}
	key := orderIDKey(order.OrderID)	// get the redis key as string
	//res := r.Client.Set(ctx, key, string(data), 0)	// cast the byte array into a string, but Set will override data that alr exists w the same key


	//	CREATE A TRANSACTION SO THEY ARE ATOMIC

	txn := r.Client.TxPipeline()	//	after creating transaction, replace all client calls to calls for transaction object

	res := txn.SetNx(ctx, key, string(data), 0)	// cast the byte array into a string
	if err := res.Err(); err != nil {
		txn.Discard()							//	if there are any errors, discard() will cancel the whole transaction
		return fmt.Error("failed to set key: %w", err)
	}

	if err := txn.SAdd(ctx,"orders",key).Err(); err != nil {	//	name of the set is "orders"
		txn.Discard()
		return fmt.Errorf("Failed to add to orders set: %w",err)
	}

	//	Final transaction will only commit when we exec() it, ensure data guarantee so no partial states
	if _, err := txn.Exec(ctx); err != nil{
		return fmt.Errorf("Failed to exec transation: %w",err)
	}

	return nil
}

//	create custom error
var ErrNotExist = errors.New("Order does not exist!")	//	from errors package


func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)

	res, err := r.Client.Get(ctx,key).Result()
	if err.Is(err, redis.Nil) {	// check if error returned is a redis.Nil error, if yes means value cannot be found from the key

		return model.Order{}, ErrNotExist	//	return empty order

	}else if err != nil {

		return model.Order{}, fmt.Errorf("Error getting order: %w", err)
	}

	var order model.Order
	err = json.Unmarshall([]byte(value), &order)	//	will write the byte array or string back into a json, which according to the struct tags, will reconstuct the order object
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode json: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	//	create transaction
	txn := r.Client.TxPipeline()	

	res, err := txn.Del(ctx,key).Result()
	if err.Is(err, redis.Nil) {	// check if error returned is a redis.Nil error, if yes means value cannot be found from the key

		txn.Discard()
		return ErrNotExist	//	return empty order error

	}else if err != nil {

		txn.Discard()
		return model.Order{}, fmt.Errorf("Error getting order: %w", err)
	}

	if err := txn.SRem(ctx,"orders",key).Err(); err != nil{		//	name of the set is "orders"

		txn.Discard()
		return fmt.Errorf("Failed to remove order from set: %w",err)
	}

	if _, err := txn.Exec(ctx); err != nil{

		return fmt.Errorf("Failed to exec transation: %w",err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error{	// can fail incase of failed communication or duplicate ID, so need to return error
	/*
	*	Note for receiver functions in Golang, must take in pointer to object else the function/method will operate on a COPY of that object
	*	However Golang automatically deferences the pointer in the function. Tricky!
	*
	*/
	data, err := json.Marshal(order)	//	encoding order struct into byte array
	if err != nil{
		return fmt.Error("Failed to encode order: %w", err)
	}
	key := orderIDKey(order.OrderID)	// get the redis key as string
	//res := r.Client.Set(ctx, key, string(data), 0)	// cast the byte array into a string, but Set will override data that alr exists w the same key
	res, err := r.Client.SetXX(ctx, key, string(data), 0)	// cast the byte array into a string, only set value when it alr exists
	if err.Is(err,redis.Nil){
		return ErrNotExist
	}else if err != nil{
		return fmt.Error("failed to set order: %w", err)
	}

	return nil
}

type FindAllPage struct {
	Size uint	// count
	Offset uint	//	cursor
}

type FindResult struct {
	Orders []model.Order	//orders set
	Cursor uint64	// the next cursor so the caller knows where to pickup from the paging process
}



func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {

	res := r.Client.SScan(ctx,"orders",page.Offset,"*"), int64(page.Size)
	keys, cur, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get order ids: %w",err)
	}

	// check if keys array is empty, if yes no point getting the results, just return empty result
	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Oder{},}, nil
	}


	xs, err := r.Client.MGet(ctx,keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get orders: %w",err)
	}

	orders := make([]model.Order,len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order

		err := json.Unmarshall([]byte(x),&order)
		if err != nil {
			return FindResult{}, fmt.Errorf("Failed to decode order json: %w",err)
		}

		orders[i] = order
	}


	return FindResult{
		Orders: orders,
		Cursor: cursor
	}, nil
}





