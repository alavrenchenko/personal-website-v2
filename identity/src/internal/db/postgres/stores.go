// Copyright 2023 Alexey Lavrenchenko. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"errors"
	"fmt"

	authenticationstores "personal-website-v2/identity/src/internal/authentication/stores"
	clientmodels "personal-website-v2/identity/src/internal/clients/models"
	clientstores "personal-website-v2/identity/src/internal/clients/stores"
	permissionstores "personal-website-v2/identity/src/internal/permissions/stores"
	rolestores "personal-website-v2/identity/src/internal/roles/stores"
	sessionmodels "personal-website-v2/identity/src/internal/sessions/models"
	sessionstores "personal-website-v2/identity/src/internal/sessions/stores"
	useragentmodels "personal-website-v2/identity/src/internal/useragents/models"
	useragentstores "personal-website-v2/identity/src/internal/useragents/stores"
	userstores "personal-website-v2/identity/src/internal/users/stores"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const (
	// identityCategory = "Identity"

	// UserStore, UserPersonalInfoStore, UserRoleAssignmentStore.
	userCategory = "User"

	// WebClientStore.
	webClientCategory = "WebClient"

	// MobileClientStore.
	mobileClientCategory = "MobileClient"

	// GroupRoleAssignmentStore.
	userGroupCategory = "UserGroup"

	// RoleStore, RolesStateStore.
	roleCategory = "Role"

	// RoleAssignmentStore.
	roleAssignmentCategory = "RoleAssignment"

	// PermissionStore, PermissionGroupStore, RolePermissionStore.
	permissionCategory = "Permission"

	// WebUserAgentStore, UserAgentWebSessionStore.
	webUserAgentCategory = "WebUserAgent"

	// MobileUserAgentStore, UserAgentMobileSessionStore.
	mobileUserAgentCategory = "MobileUserAgent"

	// UserWebSessionStore.
	userWebSessionCategory = "User_WebSession"

	// UserMobileSessionStore.
	userMobileSessionCategory = "User_MobileSession"

	// TokenEncryptionKeyStore.
	authnCategory = "Authn"
)

type Stores interface {
	UserStore() *userstores.UserStore
	UserPersonalInfoStore() *userstores.UserPersonalInfoStore
	WebClientStore() *clientstores.ClientStore
	MobileClientStore() *clientstores.ClientStore
	RoleStore() *rolestores.RoleStore
	RoleAssignmentStore() *rolestores.RoleAssignmentStore
	UserRoleAssignmentStore() *rolestores.UserRoleAssignmentStore
	GroupRoleAssignmentStore() *rolestores.GroupRoleAssignmentStore
	RolesStateStore() *rolestores.RolesStateStore
	PermissionStore() *permissionstores.PermissionStore
	PermissionGroupStore() *permissionstores.PermissionGroupStore
	RolePermissionStore() *permissionstores.RolePermissionStore
	WebUserAgentStore() *useragentstores.UserAgentStore
	MobileUserAgentStore() *useragentstores.UserAgentStore
	UserWebSessionStore() *sessionstores.UserSessionStore
	UserMobileSessionStore() *sessionstores.UserSessionStore
	UserAgentWebSessionStore() *sessionstores.UserAgentSessionStore
	UserAgentMobileSessionStore() *sessionstores.UserAgentSessionStore
	TokenEncryptionKeyStore() *authenticationstores.TokenEncryptionKeyStore
	Init(databases map[string]*postgres.Database) error
}

var _ postgres.Stores = (Stores)(nil)

type stores struct {
	userStore                   *userstores.UserStore
	userPersonalInfoStore       *userstores.UserPersonalInfoStore
	webClientStore              *clientstores.ClientStore
	mobileClientStore           *clientstores.ClientStore
	roleStore                   *rolestores.RoleStore
	roleAssignmentStore         *rolestores.RoleAssignmentStore
	userRoleAssignmentStore     *rolestores.UserRoleAssignmentStore
	groupRoleAssignmentStore    *rolestores.GroupRoleAssignmentStore
	rolesStateStore             *rolestores.RolesStateStore
	permissionStore             *permissionstores.PermissionStore
	permissionGroupStore        *permissionstores.PermissionGroupStore
	rolePermissionStore         *permissionstores.RolePermissionStore
	webUserAgentStore           *useragentstores.UserAgentStore
	mobileUserAgentStore        *useragentstores.UserAgentStore
	userWebSessionStore         *sessionstores.UserSessionStore
	userMobileSessionStore      *sessionstores.UserSessionStore
	userAgentWebSessionStore    *sessionstores.UserAgentSessionStore
	userAgentMobileSessionStore *sessionstores.UserAgentSessionStore
	tokenEncryptionKeyStore     *authenticationstores.TokenEncryptionKeyStore
	loggerFactory               logging.LoggerFactory[*context.LogEntryContext]
	isInitialized               bool
}

var _ Stores = (*stores)(nil)

func NewStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) Stores {
	return &stores{
		loggerFactory: loggerFactory,
	}
}

func (s *stores) UserStore() *userstores.UserStore {
	return s.userStore
}

func (s *stores) UserPersonalInfoStore() *userstores.UserPersonalInfoStore {
	return s.userPersonalInfoStore
}

func (s *stores) WebClientStore() *clientstores.ClientStore {
	return s.webClientStore
}

func (s *stores) MobileClientStore() *clientstores.ClientStore {
	return s.mobileClientStore
}

func (s *stores) RoleStore() *rolestores.RoleStore {
	return s.roleStore
}

func (s *stores) RoleAssignmentStore() *rolestores.RoleAssignmentStore {
	return s.roleAssignmentStore
}

func (s *stores) UserRoleAssignmentStore() *rolestores.UserRoleAssignmentStore {
	return s.userRoleAssignmentStore
}

func (s *stores) GroupRoleAssignmentStore() *rolestores.GroupRoleAssignmentStore {
	return s.groupRoleAssignmentStore
}

func (s *stores) RolesStateStore() *rolestores.RolesStateStore {
	return s.rolesStateStore
}

func (s *stores) PermissionStore() *permissionstores.PermissionStore {
	return s.permissionStore
}

func (s *stores) PermissionGroupStore() *permissionstores.PermissionGroupStore {
	return s.permissionGroupStore
}

func (s *stores) RolePermissionStore() *permissionstores.RolePermissionStore {
	return s.rolePermissionStore
}

func (s *stores) WebUserAgentStore() *useragentstores.UserAgentStore {
	return s.webUserAgentStore
}

func (s *stores) MobileUserAgentStore() *useragentstores.UserAgentStore {
	return s.mobileUserAgentStore
}

func (s *stores) UserWebSessionStore() *sessionstores.UserSessionStore {
	return s.userWebSessionStore
}

func (s *stores) UserMobileSessionStore() *sessionstores.UserSessionStore {
	return s.userMobileSessionStore
}

func (s *stores) UserAgentWebSessionStore() *sessionstores.UserAgentSessionStore {
	return s.userAgentWebSessionStore
}

func (s *stores) UserAgentMobileSessionStore() *sessionstores.UserAgentSessionStore {
	return s.userAgentMobileSessionStore
}

func (s *stores) TokenEncryptionKeyStore() *authenticationstores.TokenEncryptionKeyStore {
	return s.tokenEncryptionKeyStore
}

// databases: map[DataCategory]Database
func (s *stores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.stores.Init] stores have already been initialized")
	}

	database, ok := databases[userCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", userCategory)
	}

	userStore, err := userstores.NewUserStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new user store: %w", err)
	}

	userPersonalInfoStore, err := userstores.NewUserPersonalInfoStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new store of users' personal info: %w", err)
	}

	userRoleAssignmentStore, err := rolestores.NewUserRoleAssignmentStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new user role assignment store: %w", err)
	}

	database, ok = databases[webClientCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", webClientCategory)
	}

	webClientStore, err := clientstores.NewClientStore(clientmodels.ClientTypeWeb, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new web client store: %w", err)
	}

	database, ok = databases[mobileClientCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", mobileClientCategory)
	}

	mobileClientStore, err := clientstores.NewClientStore(clientmodels.ClientTypeMobile, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new mobile client store: %w", err)
	}

	database, ok = databases[userGroupCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", userGroupCategory)
	}

	groupRoleAssignmentStore, err := rolestores.NewGroupRoleAssignmentStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new group role assignment store: %w", err)
	}

	database, ok = databases[roleCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", roleCategory)
	}

	roleStore, err := rolestores.NewRoleStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new role store: %w", err)
	}

	rolesStateStore, err := rolestores.NewRolesStateStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new store of the state of roles: %w", err)
	}

	database, ok = databases[roleAssignmentCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", roleAssignmentCategory)
	}

	roleAssignmentStore, err := rolestores.NewRoleAssignmentStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new role assignment store: %w", err)
	}

	database, ok = databases[permissionCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", permissionCategory)
	}

	permissionStore, err := permissionstores.NewPermissionStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new permission store: %w", err)
	}

	permissionGroupStore, err := permissionstores.NewPermissionGroupStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new permission group store: %w", err)
	}

	rolePermissionStore, err := permissionstores.NewRolePermissionStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new role permission store: %w", err)
	}

	database, ok = databases[webUserAgentCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", webUserAgentCategory)
	}

	webUserAgentStore, err := useragentstores.NewUserAgentStore(useragentmodels.UserAgentTypeWeb, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new web user agent store: %w", err)
	}

	userAgentWebSessionStore, err := sessionstores.NewUserAgentSessionStore(sessionmodels.UserAgentSessionTypeWeb, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new store of web sessions of user agents: %w", err)
	}

	database, ok = databases[mobileUserAgentCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", mobileUserAgentCategory)
	}

	mobileUserAgentStore, err := useragentstores.NewUserAgentStore(useragentmodels.UserAgentTypeMobile, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new mobile user agent store: %w", err)
	}

	userAgentMobileSessionStore, err := sessionstores.NewUserAgentSessionStore(sessionmodels.UserAgentSessionTypeMobile, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new store of mobile sessions of user agents: %w", err)
	}

	database, ok = databases[userWebSessionCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", userWebSessionCategory)
	}

	userWebSessionStore, err := sessionstores.NewUserSessionStore(sessionmodels.UserSessionTypeWeb, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new store of users' web sessions: %w", err)
	}

	database, ok = databases[userMobileSessionCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", userMobileSessionCategory)
	}

	userMobileSessionStore, err := sessionstores.NewUserSessionStore(sessionmodels.UserSessionTypeMobile, database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new store of users' mobile sessions: %w", err)
	}

	database, ok = databases[authnCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", authnCategory)
	}

	tokenEncryptionKeyStore, err := authenticationstores.NewTokenEncryptionKeyStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new token encryption key store: %w", err)
	}

	s.userStore = userStore
	s.userPersonalInfoStore = userPersonalInfoStore
	s.webClientStore = webClientStore
	s.mobileClientStore = mobileClientStore
	s.roleStore = roleStore
	s.roleAssignmentStore = roleAssignmentStore
	s.userRoleAssignmentStore = userRoleAssignmentStore
	s.groupRoleAssignmentStore = groupRoleAssignmentStore
	s.rolesStateStore = rolesStateStore
	s.permissionStore = permissionStore
	s.permissionGroupStore = permissionGroupStore
	s.rolePermissionStore = rolePermissionStore
	s.webUserAgentStore = webUserAgentStore
	s.mobileUserAgentStore = mobileUserAgentStore
	s.userWebSessionStore = userWebSessionStore
	s.userMobileSessionStore = userMobileSessionStore
	s.userAgentWebSessionStore = userAgentWebSessionStore
	s.userAgentMobileSessionStore = userAgentMobileSessionStore
	s.tokenEncryptionKeyStore = tokenEncryptionKeyStore
	s.isInitialized = true
	return nil
}
