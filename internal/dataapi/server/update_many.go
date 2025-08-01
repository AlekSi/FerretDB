// Copyright 2021 FerretDB Inc.
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

package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"

	"github.com/FerretDB/wire/wirebson"

	"github.com/FerretDB/FerretDB/v2/internal/dataapi/api"
	"github.com/FerretDB/FerretDB/v2/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/v2/internal/util/must"
)

// UpdateMany implements [ServerInterface].
func (s *Server) UpdateMany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if s.l.Enabled(ctx, slog.LevelDebug) {
		s.l.DebugContext(ctx, fmt.Sprintf("Request:\n%s", must.NotFail(httputil.DumpRequest(r, true))))
	}

	var req api.UpdateRequestBody
	if err := decodeJSONRequest(r, &req); err != nil {
		http.Error(w, lazyerrors.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	updateDoc, err := prepareDocument(
		"q", req.Filter,
		"u", req.Update,
		"upsert", req.Upsert,
		"multi", true,
	)
	if err != nil {
		http.Error(w, lazyerrors.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	msg, err := prepareRequest(
		"update", req.Collection,
		"$db", req.Database,
		"updates", wirebson.MustArray(updateDoc),
	)
	if err != nil {
		http.Error(w, lazyerrors.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	resp, err := s.handler.Handle(ctx, msg)
	if err != nil {
		http.Error(w, lazyerrors.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if !resp.OK() {
		s.writeJSONError(ctx, w, resp)
		return
	}

	res := wirebson.MustDocument(
		"matchedCount", resp.Document().Get("n"),
		"modifiedCount", resp.Document().Get("nModified"),
	)

	if upsertedRaw := resp.Document().Get("upserted"); upsertedRaw != nil {
		upserted := must.NotFail(upsertedRaw.(wirebson.AnyArray).Decode())

		if upserted.Len() > 0 {
			item := must.NotFail(upserted.Get(0).(wirebson.AnyDocument).Decode())

			must.NoError(res.Add("upsertedId", item.Get("_id")))
		}
	}

	s.writeJSONResponse(ctx, w, res)
}
