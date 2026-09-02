package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/ArminDashti/knowledge-center-api/internal/models"
	"github.com/ArminDashti/knowledge-center-api/internal/store"
	"github.com/ArminDashti/knowledge-center-api/internal/worker"
	"github.com/gin-gonic/gin"
)

type API struct {
	store  *store.Store
	worker *worker.Worker
}

func New(st *store.Store, w *worker.Worker) *API {
	return &API{store: st, worker: w}
}

func (a *API) Register(r *gin.Engine) {
	r.GET("/api/health", a.Health)

	api := r.Group("/api")
	{
		api.POST("/imports", a.CreateImport)
		api.GET("/imports", a.ListImports)
		api.GET("/imports/:id", a.GetImport)
		api.GET("/imports/:id/pages/:pageId", a.GetPage)

		api.GET("/scrape-profiles", a.ListProfiles)
		api.POST("/scrape-profiles", a.CreateProfile)
		api.GET("/scrape-profiles/:id", a.GetProfile)
		api.PUT("/scrape-profiles/:id", a.UpdateProfile)
		api.DELETE("/scrape-profiles/:id", a.DeleteProfile)
	}
}

func (a *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) CreateImport(c *gin.Context) {
	var req models.CreateImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "url is required"})
		return
	}
	raw := strings.TrimSpace(req.URL)
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		if !strings.Contains(raw, "://") {
			parsed, err = url.Parse("https://" + raw)
		}
	}
	if err != nil || parsed == nil || parsed.Host == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid url"})
		return
	}
	if parsed.Scheme == "" {
		parsed.Scheme = "https"
	}

	item, err := a.store.CreateImport(c.Request.Context(), parsed.String(), strings.ToLower(parsed.Hostname()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	a.worker.Enqueue(item.ID)
	c.JSON(http.StatusAccepted, item)
}

func (a *API) ListImports(c *gin.Context) {
	items, err := a.store.ListImports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (a *API) GetImport(c *gin.Context) {
	item, err := a.store.GetImport(c.Request.Context(), c.Param("id"))
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "import not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	pages, err := a.store.ListPagesByImport(c.Request.Context(), item.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.ImportDetail{Import: *item, Pages: pages})
}

func (a *API) GetPage(c *gin.Context) {
	page, err := a.store.GetPage(c.Request.Context(), c.Param("id"), c.Param("pageId"))
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "page not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, page)
}

func (a *API) ListProfiles(c *gin.Context) {
	items, err := a.store.ListProfiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (a *API) CreateProfile(c *gin.Context) {
	var req models.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "host is required"})
		return
	}
	profile, err := a.store.CreateProfile(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, profile)
}

func (a *API) GetProfile(c *gin.Context) {
	profile, err := a.store.GetProfile(c.Request.Context(), c.Param("id"))
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (a *API) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "host is required"})
		return
	}
	profile, err := a.store.UpdateProfile(c.Request.Context(), c.Param("id"), req)
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (a *API) DeleteProfile(c *gin.Context) {
	if err := a.store.DeleteProfile(c.Request.Context(), c.Param("id")); errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "profile not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
