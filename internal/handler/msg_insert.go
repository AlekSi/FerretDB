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

package handler

import (
	"context"

	"github.com/FerretDB/wire"
	"github.com/FerretDB/wire/wirebson"
	"github.com/jackc/pgx/v5"

	"github.com/FerretDB/FerretDB/v2/internal/documentdb/documentdb_api"
	"github.com/FerretDB/FerretDB/v2/internal/handler/middleware"
	"github.com/FerretDB/FerretDB/v2/internal/mongoerrors"
	"github.com/FerretDB/FerretDB/v2/internal/util/lazyerrors"
)

// msgInsert implements `insert` command.
//
// The passed context is canceled when the client connection is closed.
func (h *Handler) msgInsert(connCtx context.Context, req *middleware.Request) (*middleware.Response, error) {
	doc, spec, seq, err := req.WireBody().(*wire.OpMsg).Sections()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	if _, _, err = h.s.CreateOrUpdateByLSID(connCtx, doc); err != nil {
		return nil, err
	}

	dbName, err := getRequiredParam[string](doc, "$db")
	if err != nil {
		return nil, err
	}

	var res wirebson.RawDocument

	err = h.p.WithConn(func(conn *pgx.Conn) error {
		res, _, err = documentdb_api.Insert(connCtx, conn, h.L, dbName, spec, seq)
		return err
	})
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return middleware.ResponseDoc(req, mongoerrors.MapWriteErrors(connCtx, res))
}
