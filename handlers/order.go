package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/generalr3d/learn-go-api-chi-redis/model"
	"github.com/generalr3d/learn-go-api-chi-redis/repository/order"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Order struct {
	Repo *order.RedisRepo //	contains the redis client

}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an order")
	var body struct { // anon struct
		CustomerID uuid.UUID        `json: "customer_id"`
		LineItems  []model.LineItem `json: "line_items"`
	}

	//create new json decoder that decodes the reuqest body
	//	passing in our anon struct basically extracts the 2 fields we want into a struct for us
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	//	create a new order for us to store into redis
	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}
	//finally insert the new order into redis repo
	err := o.Repo.Insert(r.Context(), order) //	each http request sent from client would already have a context
	if err != nil {
		fmt.Println("Failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError) // this is how to return status code to the client : 500
		return
	}

	//	if no issues
	res, err := json.Marshal(order) //	marshall means convert from json to string
	if err != nil {
		fmt.Println("failed to marshall:", err)
		w.WriteHeader(http.StatusInternalServerError) // this is how to return status code to the client : 500
		return

	}
	w.Write(res)                      //	write the marshalled json to the response
	w.WriteHeader(http.StatusCreated) //	code: 201

}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all orders")
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize) // converts string into int, if error, means the cursor provided was not numerical
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//	fetch from redis repo
	const size = 50
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("Failed to find all, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"` //	next cursor to use, if cursor is 0 returned from redis sscan means no more pages, so next field will become empoty
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal,", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return this data of type json, after encoding, back to the client
	w.Write(data)

}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get an order by ID")

	idParam := chi.URLParam(r, "id")

	const decimal = 10
	const bitSize = 64
	orderID, err := strconv.ParseUint(idParam, decimal, bitSize) // converts string into int, if error, means the cursor provided was not numerical
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theOrder, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound) // code: 404
		return
	} else if err != nil {
		fmt.Println("Failed to find by ID, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//	instead of creating a custom struct then marshall into json, this time we directly use json-encoder to encode into a json
	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("Failed to marshal, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")

	var body struct { // anon struct
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const decimal = 10
	const bitSize = 64
	orderID, err := strconv.ParseUint(idParam, decimal, bitSize) // converts string into int, if error, means the cursor provided was not numerical
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// need to pull out the existing order in order to update it and store it again
	theOrder, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound) // code: 404
		return
	} else if err != nil {
		fmt.Println("Failed to find by ID, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//	set guard conditions for which the order can be updated only if such conditions are met

	const completedStatus = "completed"
	const shippedStatus = "shipped"

	// the order is ship first, then completed
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:

		// if alr shipped and the shipped at time is alr there then throw error, else update and set the shipped time to be now
		if theOrder.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now
	case completedStatus:

		// if alr completed and the completion time is alr there OR the shiiped time is somehow still nil, throw error
		if theOrder.CompletedAt != nil || theOrder.ShippedAt == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		//	 bad request
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	// update the repo back
	err = o.Repo.Update(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to insert, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send back the updated one to the client
	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("Failed to marshal, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an order by ID")

	idParam := chi.URLParam(r, "id")

	const decimal = 10
	const bitSize = 64
	orderID, err := strconv.ParseUint(idParam, decimal, bitSize) // converts string into int, if error, means the cursor provided was not numerical
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
