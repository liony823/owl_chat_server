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

package constant

import "github.com/openimsdk/chat/pkg/protocol/constant"

const (
	// verificationCode used for.
	VerificationCodeForRegister      = 1 // Register
	VerificationCodeForResetPassword = 2 // Reset password
	VerificationCodeForLogin         = 3 // Login
)

const LogFileName = "chat.log"

// block unblock.
const (
	BlockUser   = 1
	UnblockUser = 2
)

// AccountType.
const (
	Email   = "email"
	Phone   = "phone"
	Account = "account"
)

// Mode.
const (
	UserMode  = "user"
	AdminMode = "admin"
)

const DefaultAdminLevel = 100

// user level.
const (
	NormalAdmin       = 80
	AdvancedUserLevel = 100
)

// AddFriendCtrl.
const (
	OrdinaryUserAddFriendEnable  = 1  // Allow ordinary users to add friends
	OrdinaryUserAddFriendDisable = -1 // Do not allow ordinary users to add friends
)

const (
	NormalUser = 1
	AdminUser  = 2
)

// mini-app
const (
	StatusOnShelf = 1 // OnShelf
	StatusUnShelf = 2 // UnShelf
)

const (
	LimitNil             = 0 // None
	LimitEmpty           = 1 // Neither are restricted
	LimitOnlyLoginIP     = 2 // Only login is restricted
	LimitOnlyRegisterIP  = 3 // Only registration is restricted
	LimitLoginIP         = 4 // Restrict login
	LimitRegisterIP      = 5 // Restrict registration
	LimitLoginRegisterIP = 6 // Restrict both login and registration
)

const (
	InvitationCodeAll    = 0 // All
	InvitationCodeUsed   = 1 // Used
	InvitationCodeUnused = 2 // Unused
)

const (
	RpcOpUserID   = constant.OpUserID
	RpcOpUserType = "opUserType"
)

const RpcCustomHeader = constant.RpcCustomHeader

const NeedInvitationCodeRegisterConfigKey = "needInvitationCodeRegister"

const (
	DefaultAllowVibration = 1
	DefaultAllowBeep      = 1
	DefaultAllowAddFriend = 1
)

const (
	FinDAllUser    = 0
	FindNormalUser = 1
)

const CtxApiToken = "api-token"

const (
	EmailRegister = 1
	PhoneRegister = 2
)

const (
	Follow    = 0
	Subscribe = 1
	Reply     = 2
	Like      = 3
	Collect   = 4
)

const (
	NotAllow = 0
	Allow    = 1
)

const (
	NotLiked = 0
	Liked    = 1
)

const (
	NotCollected = 0
	Collected    = 1
)

const (
	NotForwarded = 0
	Forwarded    = 1
)

const (
	NotCommented = 0
	Commented    = 1
)

const (
	NotFollowed = 0
	Followed    = 1
)

const (
	NotSubscribed = 0
	Subscribed    = 1
)

const (
	PostMediaTypePicture = 0
	PostMediaTypeVideo   = 1
)

const (
	Pinned   = 1
	UnPinned = 0
)
