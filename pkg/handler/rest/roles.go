package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oeoen/policy/helper"
	"github.com/oeoen/policy/pkg/police"
)

func UpsertPolicyRole(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := &police.RBAC{}
		if err := c.Bind(body); err != nil {
			cn.Logger().Error(err)
			return err
		}
		tenant := c.Param("tenant")
		body.Tenant = tenant
		if err := m.UpsertRole(c.Request().Context(), body); err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusCreated, body)
	}
}

func GetPolicyRole(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant, policy := c.Param("tenant"), c.Param("policy")
		acl, err := m.GetRolePolicy(c.Request().Context(), tenant, policy)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, acl)
	}
}

func GetSubjectRoles(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant, subject := c.Param("tenant"), c.Param("subject")
		acl, err := m.GetSubjectRoles(c.Request().Context(), tenant, subject)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, acl)
	}
}

func GetRoleSubjects(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant, policy := c.Param("tenant"), c.Param("policy")
		acl, err := m.GetRoleSubjects(c.Request().Context(), tenant, policy)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusOK, acl)
	}
}

func DeleteSubjectsRole(m police.Manager, cn police.Configuration) echo.HandlerFunc {
	return func(c echo.Context) error {
		tenant, policy, subject := c.Param("tenant"), c.Param("policy"), c.Param("subject")
		err := m.DeleteRole(c.Request().Context(), tenant, subject, policy)
		if err != nil {
			cn.Logger().Error(err)
			return responseError(c, err)
		}
		return c.JSON(http.StatusAccepted, helper.Response{Message: "Deleted"})
	}
}
