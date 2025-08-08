package common

import (
	"net/http"
	"testYTask/internal/domain/models"

	"github.com/gin-gonic/gin"
)

// ResponseBadRequest отправляет ответ с кодом 400 (Bad Request).
// Используется, когда клиент прислал некорректные данные.
//
// Параметры:
//   - c: контекст Gin, через который формируется ответ
//   - msg: сообщение об ошибке, которое будет возвращено в поле "message"
func ResponseBadRequest(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
		Status:  statusError,
		Message: msg,
	})
}

// ResponseServerError отправляет ответ с кодом 500 (Internal Server Error)
// Используется при ошибках на стороне сервера.
//
// Параметры:
//   - c: контекст Gin, через который формируется ответ
//   - msg: сообщение об ошибке для пользователя или разработчика
func ResponseServerError(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
		Status:  statusError,
		Message: msg,
	})
}

// ResponseNotFound отправляет ответ с кодом 404 (Not Found).
// Используется, когда запрошенный ресурс не найден.
//
// Параметры:
//   - c: контекст Gin, через который формируется ответ
//   - msg: сообщение, описывающее, что именно не найдено
func ResponseNotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, models.Response{
		Status:  statusNotFound,
		Message: msg,
	})
}

// ResponseSuccess отправляет успешный ответ с кодом 200 (OK).
// Используется для возврата успешного результата.
//
// Параметры:
//   - c: контекст Gin, через который формируется ответ
//   - msg: сообщение об успешном выполнении операции
//   - data: любые дополнительные данные, которые будут переданы в поле "data"
func ResponseSuccess(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, models.Response{
		Status:  statusOK,
		Message: msg,
		Data:    data,
	})
}
