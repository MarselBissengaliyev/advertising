package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	error_message "github.com/MarselBissengaliyev/advertising/pkg/error_messaage"
	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/labstack/echo/v4"
)

type getAllListsResponse struct {
	Data []model.GetAllAdvertsResponse `json:"data"`
}

// @Summary Get all adverts
// @Description Get all advertising items
// @Tags adverts
// @Produce json
// @Success 200 {object} getAllListsResponse
// @Failure default {object} errorResponse
// @Param sort query string false "Sort" Format(string)
// @Param page query string false "Page" Format(int)
// @Router /api/adverts/ [get]
func (h *Handler) getAllAdverts(c echo.Context) error {
	limit := 10

	orderBy := c.QueryParam("sort")
	validOrderByValues := [4]string{"price_desc", "price_asc", "created_at_desc", "created_at_asc"}
	orderByIsValid := false
	for _, value := range validOrderByValues {
		if orderBy == value {
			orderByIsValid = true
			break
		}
	}
	if !orderByIsValid {
		orderBy = "price_asc"
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}

	adverts, err := h.services.Advert.GetAll(limit, page, orderBy)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	err = c.JSON(http.StatusOK, getAllListsResponse{
		Data: adverts,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	
	return nil
}

// @Summary Get advert by id
// @Description Get advert by id
// @Tags adverts
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure default {object} errorResponse
// @Param fields query string false "Fields to include (comma-separated)" Format(string)
// @Param id path string false "Advert ID" Format(int)
// @Router /api/adverts/{id} [get]
func (h *Handler) getAdvertById(c echo.Context) error {
	var photos *[]string
	var response map[string]interface{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, error_message.InvalidIdParam)
		return err
	}

	fields := c.QueryParam("fields")
	parsedFields := strings.Split(fields, ",")
	fmt.Println(parsedFields)

	advert, err := h.services.Advert.GetById(id, parsedFields)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	response = map[string]interface{}{
		"id":         advert.Id,
		"title":      advert.Title,
		"price":      advert.Price,
		"main_photo": advert.MainPhoto,
	}

	if advert.Description != "" {
		response["description"] = advert.Description
	}

	if len(advert.Photos) > 0 {
		err = json.Unmarshal([]byte(advert.Photos), &photos)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return err
		}

		response["photos"] = photos
	}

	err = c.JSON(http.StatusOK, response)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	return nil
}

// @Summary Create a new advert
// @Description Create a new advertising item
// @Tags adverts
// @Accept json
// @Produce json
// @Param input body model.CreateAdvertsBody true "Advert details to create"
// @Success 201 {object} map[string]interface{} "Created"
// @Router /api/adverts/ [post]
func (h *Handler) createAdvert(c echo.Context) error {
	var input model.CreateAdvertsBody
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := h.services.Advert.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	err = c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	});
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}
