// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package imapi

import (
	"github.com/openimsdk/chat/pkg/protocol/auth"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/chat/pkg/protocol/friend"
	"github.com/openimsdk/chat/pkg/protocol/group"
	"github.com/openimsdk/chat/pkg/protocol/user"
	"github.com/openimsdk/protocol/msggateway"
)

// im caller.
var (
	importFriend        = NewApiCaller[friend.ImportFriendReq, friend.ImportFriendResp]("/friend/import_friend")
	userToken           = NewApiCaller[auth.UserTokenReq, auth.UserTokenResp]("/auth/user_token")
	parseToken          = NewApiCaller[auth.ParseTokenReq, auth.ParseTokenResp]("/auth/parse_token")
	inviteToGroup       = NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")
	updateUserInfo      = NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")
	registerUser        = NewApiCaller[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register")
	forceOffLine        = NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp]("/auth/force_logout")
	getGroupsInfo       = NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info")
	registerUserCount   = NewApiCaller[user.UserRegisterCountReq, user.UserRegisterCountResp]("/statistics/user/register")
	friendUserIDs       = NewApiCaller[friend.GetFriendIDsReq, friend.GetFriendIDsResp]("/friend/get_friend_id")
	accountCheck        = NewApiCaller[user.AccountCheckReq, user.AccountCheckResp]("/user/account_check")
	allUserOnlineStatus = NewApiCaller[msggateway.GetUsersOnlineStatusReq, []msggateway.GetUsersOnlineStatusResp_SuccessResult]("/user/get_users_online_status")
	usersOnlineTime     = NewApiCaller[chat.GetUsersTimeReq, chat.GetUsersTimeResp]("/user/get_users_time")
)
