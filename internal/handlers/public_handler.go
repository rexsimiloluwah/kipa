package handlers

import (
	"keeper/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IPublicRoutesHandler interface {
	HealthCheck(c echo.Context) error
	GetAPIKeyPermissionsList(c echo.Context) error
	GetBucketPermissionsList(c echo.Context) error
}

type PublicRoutesHandler struct {
}

func NewPublicRoutesHandler() IPublicRoutesHandler {
	return &PublicRoutesHandler{}
}

// HealthCheck godoc
// @Summary      HealthCheck
// @Description  Check if the server is healthy
// @Tags         PublicRoutes
// @Produce      json
// @Success      200  {object} 	models.SuccessResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /public/healthcheck [get]
func (h *PublicRoutesHandler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Server is healthy!")
}

// GetAPIKeyPermissionsList godoc
// @Summary      GetAPIKeyPermissionsList
// @Description  Returns a list of the API key permissions
// @Tags         PublicRoutes
// @Produce      json
// @Success      200  {object} 	models.SuccessResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /public/apikey-permissions [get]
func (h *PublicRoutesHandler) GetAPIKeyPermissionsList(c echo.Context) error {
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched api key permissions",
		Data:    models.APIKEY_PERMISSIONS,
	})
}

// GetBucketPermissionsList godoc
// @Summary      GetBucketPermissionsList
// @Description  Returns a list of the bucket permissions
// @Tags         PublicRoutes
// @Produce      json
// @Success      200  {object} 	models.SuccessResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /public/bucket-permissions [get]
func (h *PublicRoutesHandler) GetBucketPermissionsList(c echo.Context) error {
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched bucket permissions",
		Data:    models.BUCKET_PERMISSIONS,
	})
}
