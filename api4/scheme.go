// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"
)

func (api *API) InitScheme() {
	api.BaseRoutes.Schemes.Handle("", api.ApiSessionRequired(createScheme)).Methods("POST")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}", api.ApiSessionRequired(deleteScheme)).Methods("DELETE")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}", api.ApiSessionRequiredTrustRequester(getScheme)).Methods("GET")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}/patch", api.ApiSessionRequired(patchScheme)).Methods("PUT")
}

func createScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	scheme := model.SchemeFromJson(r.Body)
	if scheme == nil {
		c.SetInvalidParam("scheme")
		return
	}

	if c.App.License() == nil {
		c.Err = model.NewAppError("Api4.CreateScheme", "api.scheme.create_scheme.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	var err *model.AppError
	if scheme, err = c.App.CreateScheme(scheme); err != nil {
		c.Err = err
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(scheme.ToJson()))
	}
}

func getScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	if scheme, err := c.App.GetScheme(c.Params.SchemeId); err != nil {
		c.Err = err
		return
	} else {
		w.Write([]byte(scheme.ToJson()))
	}
}

func patchScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	patch := model.SchemePatchFromJson(r.Body)
	if patch == nil {
		c.SetInvalidParam("scheme")
		return
	}

	if c.App.License() == nil {
		c.Err = model.NewAppError("Api4.PatchScheme", "api.scheme.patch_scheme.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	scheme, err := c.App.GetScheme(c.Params.SchemeId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	if scheme, err = c.App.PatchScheme(scheme, patch); err != nil {
		c.Err = err
		return
	} else {
		c.LogAudit("")
		w.Write([]byte(scheme.ToJson()))
	}
}

func deleteScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	if c.App.License() == nil {
		c.Err = model.NewAppError("Api4.DeleteScheme", "api.scheme.delete_scheme.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	if _, err := c.App.DeleteScheme(c.Params.SchemeId); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}
