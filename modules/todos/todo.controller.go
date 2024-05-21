package todos

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nebnhoj/strand/helpers"
)

func CreateTodo(c *fiber.Ctx) error {
	var todo Todo
	//validate the request body
	if err := c.BodyParser(&todo); err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	if errs := helpers.Validator(todo); len(errs) > 0 && errs[0].Error {
		return helpers.ResponseError(c, http.StatusBadRequest, errs)
	}
	newTodo := Todo{
		Id:      uuid.NewString(),
		Name:    todo.Name,
		Details: todo.Details,
	}

	result, err := create(newTodo)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	oldRedisData, err := rdb.Get(context.Background(), c.Path()).Result()
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	var retrievedClaims []Todo
	err = json.Unmarshal([]byte(oldRedisData), &retrievedClaims)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())

	}
	retrievedClaims = append(retrievedClaims, newTodo)
	jsonData, err := json.Marshal(retrievedClaims)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())

	}
	err = rdb.Set(context.Background(), c.Path(), jsonData, 60*time.Hour).Err()
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.ResponseSuccess(c, http.StatusOK, result)
}

func FindAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	q := c.Query("q")
	data, err := rdb.Get(context.Background(), c.OriginalURL()).Result()
	if err == nil {
		var retrievedClaims []Todo
		err = json.Unmarshal([]byte(data), &retrievedClaims)
		if err != nil {
			return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return helpers.ResponseSuccess(c, http.StatusOK, retrievedClaims)
	}

	// Convert pagination parameters to integers

	todos, err := getAllTodos(q, page, limit)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	jsonData, err := json.Marshal(todos)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}
	err = rdb.Set(context.Background(), c.OriginalURL(), jsonData, 60*time.Hour).Err()
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.ResponseSuccess(c, http.StatusOK, todos)
}

type Message struct {
	Content int64 `json:"content"`
}

type Status struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}
