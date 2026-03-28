package http

import (
	"context"
	stdhttp "net/http"

	"portlyn/internal/domain"
)

func (s *Server) handleListGroups(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	items, err := s.groups.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	response := make([]map[string]any, 0, len(items))
	for _, item := range items {
		memberCount, err := s.groups.CountMembers(r.Context(), item.ID)
		if err != nil {
			s.internalError(w, err)
			return
		}
		response = append(response, map[string]any{
			"id":              item.ID,
			"name":            item.Name,
			"description":     item.Description,
			"is_system_group": item.IsSystemGroup,
			"member_count":    memberCount,
			"created_at":      item.CreatedAt,
			"updated_at":      item.UpdatedAt,
		})
	}
	writeJSON(w, stdhttp.StatusOK, response)
}

func (s *Server) handleGetGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := s.groups.GetByID(r.Context(), id)
	if err != nil {
		s.handleStoreError(w, err)
		return
	}
	members, err := s.groups.ListMembers(r.Context(), item.ID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, map[string]any{
		"id":              item.ID,
		"name":            item.Name,
		"description":     item.Description,
		"is_system_group": item.IsSystemGroup,
		"members":         members,
		"created_at":      item.CreatedAt,
		"updated_at":      item.UpdatedAt,
	})
}

func (s *Server) handleCreateGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req createGroupRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	item := &domain.Group{Name: req.Name, Description: req.Description}
	if err := s.groups.Create(r.Context(), item); err != nil {
		s.internalError(w, err)
		return
	}
	s.auth.InvalidateAll()
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "create", "group", &item.ID, item)
	writeJSON(w, stdhttp.StatusCreated, item)
}

func (s *Server) handleUpdateGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := s.groups.GetByID(r.Context(), id)
	if err != nil {
		s.handleStoreError(w, err)
		return
	}
	var req updateGroupRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if err := s.groups.Update(r.Context(), item); err != nil {
		s.internalError(w, err)
		return
	}
	s.auth.InvalidateAll()
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "update", "group", &item.ID, item)
	writeJSON(w, stdhttp.StatusOK, item)
}

func (s *Server) handleDeleteGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := s.groups.GetByID(r.Context(), id)
	if err != nil {
		s.handleStoreError(w, err)
		return
	}
	if item.IsSystemGroup {
		writeError(w, stdhttp.StatusConflict, "system_group_protected", "system groups cannot be deleted")
		return
	}
	if err := s.groups.Delete(r.Context(), id); err != nil {
		s.handleStoreError(w, err)
		return
	}
	s.auth.InvalidateAll()
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "delete", "group", &id, map[string]any{"id": id, "name": item.Name})
	w.WriteHeader(stdhttp.StatusNoContent)
}

func (s *Server) handleAddGroupMember(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	groupID, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	var req groupMemberRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	if _, err := s.users.GetByID(r.Context(), req.UserID); err != nil {
		s.handleStoreError(w, err)
		return
	}
	if err := s.groups.AddMember(r.Context(), groupID, req.UserID); err != nil {
		s.internalError(w, err)
		return
	}
	s.auth.InvalidateUser(req.UserID)
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "membership_add", "group", &groupID, map[string]any{"user_id": req.UserID})
	s.handleGetGroup(w, r)
}

func (s *Server) handleDeleteGroupMember(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	groupID, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	userID, ok := s.parseIDParam(w, r, "userId")
	if !ok {
		return
	}
	if err := s.groups.RemoveMember(r.Context(), groupID, userID); err != nil {
		s.handleStoreError(w, err)
		return
	}
	s.auth.InvalidateUser(userID)
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "membership_remove", "group", &groupID, map[string]any{"user_id": userID})
	w.WriteHeader(stdhttp.StatusNoContent)
}

func (s *Server) handleListServiceGroups(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	items, err := s.serviceGroups.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	response := make([]map[string]any, 0, len(items))
	for _, item := range items {
		response = append(response, serviceGroupResponse(item))
	}
	writeJSON(w, stdhttp.StatusOK, response)
}

func (s *Server) handleGetServiceGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := s.serviceGroups.GetByID(r.Context(), id)
	if err != nil {
		s.handleStoreError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, serviceGroupResponse(*item))
}

func (s *Server) handleCreateServiceGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req createServiceGroupRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	item := &domain.ServiceGroup{
		Name:        req.Name,
		Description: req.Description,
		DefaultAccessPolicy: domain.AccessPolicy{
			AccessMode:           req.DefaultAccessPolicy.AccessMode,
			AllowedRoles:         normalizeStringList(req.DefaultAccessPolicy.AllowedRoles),
			AllowedGroups:        domain.JSONUintSlice(req.DefaultAccessPolicy.AllowedGroups),
			AllowedServiceGroups: domain.JSONUintSlice(req.DefaultAccessPolicy.AllowedServiceGroups),
		},
		AccessMethod:       normalizeOptionalAccessMethod(req.AccessMethod),
		AccessMethodConfig: buildAccessMethodConfig(req.AccessMethod, req.AccessMethodConfig, nil),
	}
	if err := s.serviceGroups.Create(r.Context(), item); err != nil {
		s.internalError(w, err)
		return
	}
	if err := s.serviceGroups.ReplaceServices(r.Context(), item.ID, req.ServiceIDs); err != nil {
		s.internalError(w, err)
		return
	}
	created, err := s.serviceGroups.GetByID(r.Context(), item.ID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	if err := s.refreshServices(r.Context(), serviceIDsFromItems(created.Services)); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "create", "service_group", &created.ID, created)
	writeJSON(w, stdhttp.StatusCreated, serviceGroupResponse(*created))
}

func (s *Server) handleUpdateServiceGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := s.serviceGroups.GetByID(r.Context(), id)
	if err != nil {
		s.handleStoreError(w, err)
		return
	}
	affectedServiceIDs := serviceIDsFromItems(item.Services)
	var req updateServiceGroupRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.DefaultAccessPolicy != nil {
		item.DefaultAccessPolicy = domain.AccessPolicy{
			AccessMode:           req.DefaultAccessPolicy.AccessMode,
			AllowedRoles:         normalizeStringList(req.DefaultAccessPolicy.AllowedRoles),
			AllowedGroups:        domain.JSONUintSlice(req.DefaultAccessPolicy.AllowedGroups),
			AllowedServiceGroups: domain.JSONUintSlice(req.DefaultAccessPolicy.AllowedServiceGroups),
		}
	}
	if req.AccessMethod != nil {
		item.AccessMethod = normalizeOptionalAccessMethod(*req.AccessMethod)
	}
	if req.AccessMethodConfig != nil || req.AccessMethod != nil {
		method := item.AccessMethod
		if req.AccessMethod != nil {
			method = *req.AccessMethod
		}
		item.AccessMethodConfig = buildAccessMethodConfig(method, derefAccessMethodConfig(req.AccessMethodConfig), item.AccessMethodConfig)
	}
	if err := s.serviceGroups.Update(r.Context(), item); err != nil {
		s.internalError(w, err)
		return
	}
	if req.ServiceIDs != nil {
		if err := s.serviceGroups.ReplaceServices(r.Context(), item.ID, *req.ServiceIDs); err != nil {
			s.internalError(w, err)
			return
		}
	}
	updated, err := s.serviceGroups.GetByID(r.Context(), item.ID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	affectedServiceIDs = append(affectedServiceIDs, serviceIDsFromItems(updated.Services)...)
	if err := s.refreshServices(r.Context(), affectedServiceIDs); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "update", "service_group", &updated.ID, updated)
	writeJSON(w, stdhttp.StatusOK, serviceGroupResponse(*updated))
}

func (s *Server) handleDeleteServiceGroup(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := s.serviceGroups.GetByID(r.Context(), id)
	if err != nil {
		s.handleStoreError(w, err)
		return
	}
	affectedServiceIDs := serviceIDsFromItems(item.Services)
	if err := s.serviceGroups.Delete(r.Context(), id); err != nil {
		s.handleStoreError(w, err)
		return
	}
	if err := s.refreshServices(r.Context(), affectedServiceIDs); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "delete", "service_group", &id, map[string]any{"id": id, "name": item.Name})
	w.WriteHeader(stdhttp.StatusNoContent)
}

func (s *Server) handleAddServiceGroupService(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	var req serviceGroupServiceRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	if _, err := s.services.GetByID(r.Context(), req.ServiceID); err != nil {
		s.handleStoreError(w, err)
		return
	}
	if err := s.serviceGroups.AddService(r.Context(), id, req.ServiceID); err != nil {
		s.internalError(w, err)
		return
	}
	updated, err := s.serviceGroups.GetByID(r.Context(), id)
	if err != nil {
		s.internalError(w, err)
		return
	}
	if err := s.refreshServices(r.Context(), []uint{req.ServiceID}); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "membership_add", "service_group", &id, map[string]any{"service_id": req.ServiceID})
	writeJSON(w, stdhttp.StatusOK, serviceGroupResponse(*updated))
}

func (s *Server) handleDeleteServiceGroupService(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	serviceID, ok := s.parseIDParam(w, r, "serviceId")
	if !ok {
		return
	}
	if err := s.serviceGroups.RemoveService(r.Context(), id, serviceID); err != nil {
		s.internalError(w, err)
		return
	}
	if err := s.refreshServices(r.Context(), []uint{serviceID}); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "membership_remove", "service_group", &id, map[string]any{"service_id": serviceID})
	w.WriteHeader(stdhttp.StatusNoContent)
}

func (s *Server) refreshServices(ctx context.Context, serviceIDs []uint) error {
	seen := make(map[uint]struct{}, len(serviceIDs))
	for _, serviceID := range serviceIDs {
		if serviceID == 0 {
			continue
		}
		if _, ok := seen[serviceID]; ok {
			continue
		}
		seen[serviceID] = struct{}{}
		if _, err := s.proxy.ApplyServiceChange(ctx, serviceID); err != nil {
			return err
		}
	}
	return nil
}

func serviceIDsFromItems(items []domain.Service) []uint {
	ids := make([]uint, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}
