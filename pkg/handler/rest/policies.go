package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oeoen/policy/helper"
	"github.com/oeoen/policy/pkg/police"
)

func CreatePolicy(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.ACL{}
		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		if err := m.UpsertPolicy(c.Request().Context(), body); err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusCreated, body)
	}
}

func UpdatePolicy(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.ACL{}
		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		pID := parseParam(c, "police_id")
		if err := m.UpdatePolicy(c.Request().Context(), pID, body); err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusCreated, body)
	}
}

func GetPolicy(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.ACL{}
		id := parseParam(c, "policy_id")
		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		acls, err := m.GetPolicy(c.Request().Context(), id)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, acls)
	}
}

func GetPolicyTenant(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.ACL{}
		tenant := parseParam(c, "tenant")
		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		acls, err := m.FetchPolicy(c.Request().Context(), [3]string{"tenant", "=", tenant})
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, acls)
	}
}

func GetPolicyResources(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		r, err := m.GetResources(c.Request().Context())
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, r)
	}
}

func GetPolicySubjects(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		r, err := m.GetPolicySubjects(c.Request().Context())
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, r)
	}
}

func FetchPolicy(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.ACL{}
		queries := getQueries(c.QueryParams())
		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		acls, err := m.FetchPolicy(c.Request().Context(), queries...)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, acls)
	}
}

func DeletePolicy(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.ACL{}
		id := parseParam(c, "policy_id")
		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		err := m.DeletePolicy(c.Request().Context(), id)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusAccepted, helper.Response{Message: "Deleted"})
	}
}

func Enforce(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.ACL{}

		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		tenant := parseParam(c, "tenant")
		acl, err := m.Enforce(c.Request().Context(), tenant, body.Subject, body.Action, body.Resource)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		if acl.Check() != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(200, map[string]bool{"allowed": true})
	}
}
