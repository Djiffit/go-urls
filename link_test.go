package shortener_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	shortener "github.com/Djiffit/url-shortener"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var testLink = `{"id":"test","target":"helsinki.fi"}`

func TestLinks(t *testing.T) {

	t.Run("Link is properly created", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testLink))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := &shortener.LinkModel{shortener.MemoryLinkStore{map[string]string{}}}

		if assert.NoError(t, h.PostLink(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, testLink, strings.TrimSpace(rec.Body.String()))
		}

	})

	t.Run("Created link can be retrieved", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testLink))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := &shortener.LinkModel{shortener.MemoryLinkStore{map[string]string{}}}

		if assert.NoError(t, h.PostLink(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, testLink, strings.TrimSpace(rec.Body.String()))

			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(testLink))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("test")
			err := h.GetLink(c)

			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.True(t, strings.Contains(rec.Body.String(), "helsinki.fi"))
			}
		}

	})

	t.Run("Can't create a link that already exists", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testLink))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := &shortener.LinkModel{shortener.MemoryLinkStore{map[string]string{"test": "best"}}}
		err := h.PostLink(c)

		assert.EqualError(t, err, echo.NewHTTPError(http.StatusBadRequest, shortener.ErrIDExists).Error())
	})

	t.Run("Can delete a link", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(testLink))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("test")
		h := &shortener.LinkModel{shortener.MemoryLinkStore{map[string]string{"test": "best"}}}
		err := h.DeleteLink(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)

			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(testLink))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("test")
			err := h.GetLink(c)

			if err == nil {
				t.Errorf("Expected error")
			}
		}
	})

	t.Run("Get returns the desired url", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(testLink))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("test")
		h := &shortener.LinkModel{shortener.MemoryLinkStore{map[string]string{"test": "best"}}}
		err := h.GetLink(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.True(t, strings.Contains(rec.Body.String(), "best"))
		}
	})
}
