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

package chat

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mw/specialerror"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func generateNonce(size int) (string, error) {
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil
}

func validateSignature(publicKey, nonce, signature string) (b bool, err error) {
	pubKey, err := hexutil.Decode(publicKey)
	if err != nil {
		return false, err
	}
	message, err := hexutil.Decode(nonce)
	if err != nil {
		return false, err
	}
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}
	b = crypto.VerifySignature(pubKey, message, sig)
	return b, nil
}

func IsNotFound(err error) bool {
	return errs.ErrRecordNotFound.Is(specialerror.ErrCode(errs.Unwrap(err)))
}
