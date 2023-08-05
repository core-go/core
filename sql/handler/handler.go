package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	q "github.com/core-go/core/sql"
)

type Handler struct {
	DB        *sql.DB
	Transform func(s string) string
	Cache     q.TxCache
	Generate  func(ctx context.Context) (string, error)
	Error     func(context.Context, string, ...map[string]interface{})
}

const d = 120 * time.Second
func NewHandler(db *sql.DB, transform func(s string) string, cache q.TxCache, generate func(context.Context) (string, error), options... func(context.Context, string, ...map[string]interface{})) *Handler {
	var logError func(context.Context, string, ...map[string]interface{})
	if len(options) >= 1 {
		logError = options[0]
	}
	return &Handler{DB: db, Transform: transform, Cache: cache, Generate: generate, Error: logError}
}
func (h *Handler) BeginTransaction(w http.ResponseWriter, r *http.Request) {
	id, er0 := h.Generate(r.Context())
	if er0 != nil {
		http.Error(w, er0.Error(), http.StatusInternalServerError)
		return
	}
	tx, er1 := h.DB.Begin()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	ps := r.URL.Query()
	t := d
	st := ps.Get("timeout")
	if len(st) > 0 {
		i, er2 := strconv.ParseInt(st, 10, 64)
		if er2 == nil && i > 0 {
			t = time.Duration(i) * time.Second
		}
	}
	h.Cache.Put(id, tx, t)
	respond(w, http.StatusOK, id)
}
func (h *Handler) EndTransaction(w http.ResponseWriter, r *http.Request) {
	ps := r.URL.Query()
	stx := ps.Get("tx")
	if len(stx) == 0 {
		http.Error(w, "tx is required", http.StatusBadRequest)
		return
	}
	tx, er0 := h.Cache.Get(stx)
	if er0 != nil {
		http.Error(w, er0.Error(), http.StatusInternalServerError)
		return
	}
	if tx == nil {
		http.Error(w, "cannot get tx from cache. Maybe tx got timeout", http.StatusInternalServerError)
		return
	}
	rollback := ps.Get("rollback")
	if rollback == "true" {
		er1 := tx.Rollback()
		if er1 != nil {
			http.Error(w, er1.Error(), http.StatusInternalServerError)
		} else {
			h.Cache.Remove(stx)
			respond(w, http.StatusOK, "true")
		}
	} else {
		er1 := tx.Commit()
		if er1 != nil {
			http.Error(w, er1.Error(), http.StatusInternalServerError)
		} else {
			h.Cache.Remove(stx)
			respond(w, http.StatusOK, "true")
		}
	}
}
func (h *Handler) Exec(w http.ResponseWriter, r *http.Request) {
	s := q.JStatement{}
	er0 := json.NewDecoder(r.Body).Decode(&s)
	if er0 != nil {
		http.Error(w, er0.Error(), http.StatusBadRequest)
		return
	}
	s.Params = q.ParseDates(s.Params, s.Dates)
	ps := r.URL.Query()
	stx := ps.Get("tx")
	if len(stx) == 0 {
		res, er1 := h.DB.Exec(s.Query, s.Params...)
		if er1 != nil {
			handleError(w, r, 500, er1.Error(), h.Error, er1)
			return
		}
		a2, er2 := res.RowsAffected()
		if er2 != nil {
			handleError(w, r, http.StatusInternalServerError, er2.Error(), h.Error, er2)
			return
		}
		respond(w, http.StatusOK, a2)
	} else {
		tx, er0 := h.Cache.Get(stx)
		if er0 != nil {
			http.Error(w, er0.Error(), http.StatusInternalServerError)
			return
		}
		if tx == nil {
			http.Error(w, "cannot get tx from cache. Maybe tx got timeout", http.StatusInternalServerError)
			return
		}
		res, er1 := tx.Exec(s.Query, s.Params...)
		if er1 != nil {
			tx.Rollback()
			h.Cache.Remove(stx)
			handleError(w, r, 500, er1.Error(), h.Error, er1)
			return
		}
		a2, er2 := res.RowsAffected()
		if er2 != nil {
			handleError(w, r, http.StatusInternalServerError, er2.Error(), h.Error, er2)
			return
		}
		commit := ps.Get("commit")
		if commit == "true" {
			er3 := tx.Commit()
			if er3 != nil {
				handleError(w, r, http.StatusInternalServerError, er3.Error(), h.Error, er3)
				return
			}
			h.Cache.Remove(stx)
		}
		respond(w, http.StatusOK, a2)
	}
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	s := q.JStatement{}
	er0 := json.NewDecoder(r.Body).Decode(&s)
	if er0 != nil {
		http.Error(w, er0.Error(), http.StatusBadRequest)
		return
	}
	s.Params = q.ParseDates(s.Params, s.Dates)
	ps := r.URL.Query()
	stx := ps.Get("tx")
	if len(stx) == 0 {
		res, er1 := q.QueryMap(r.Context(), h.DB, h.Transform, s.Query, s.Params...)
		if er1 != nil {
			handleError(w, r, 500, er1.Error(), h.Error, er1)
			return
		}
		respond(w, http.StatusOK, res)
	} else {
		tx, er0 := h.Cache.Get(stx)
		if er0 != nil {
			http.Error(w, er0.Error(), http.StatusInternalServerError)
			return
		}
		if tx == nil {
			http.Error(w, "cannot get tx from cache. Maybe tx got timeout", http.StatusInternalServerError)
			return
		}
		res, er1 := q.QueryMapWithTx(r.Context(), tx, h.Transform, s.Query, s.Params...)
		if er1 != nil {
			handleError(w, r, 500, er1.Error(), h.Error, er1)
			return
		}
		commit := ps.Get("commit")
		if commit == "true" {
			er3 := tx.Commit()
			if er3 != nil {
				handleError(w, r, http.StatusInternalServerError, er3.Error(), h.Error, er3)
				return
			}
			h.Cache.Remove(stx)
		}
		respond(w, http.StatusOK, res)
	}
}
func (h *Handler) QueryOne(w http.ResponseWriter, r *http.Request) {
	s := q.JStatement{}
	er0 := json.NewDecoder(r.Body).Decode(&s)
	if er0 != nil {
		http.Error(w, er0.Error(), http.StatusBadRequest)
		return
	}
	s.Params = q.ParseDates(s.Params, s.Dates)
	ps := r.URL.Query()
	stx := ps.Get("tx")
	if len(stx) == 0 {
		res, er1 := q.QueryMap(r.Context(), h.DB, h.Transform, s.Query, s.Params...)
		if er1 != nil {
			handleError(w, r, 500, er1.Error(), h.Error, er1)
			return
		}
		if len(res) > 0 {
			respond(w, http.StatusOK, res[0])
			return
		}
		respond(w, http.StatusOK, nil)
	} else {
		tx, er0 := h.Cache.Get(stx)
		if er0 != nil {
			http.Error(w, er0.Error(), http.StatusInternalServerError)
			return
		}
		if tx == nil {
			http.Error(w, "cannot get tx from cache. Maybe tx got timeout", http.StatusInternalServerError)
			return
		}
		res, er1 := q.QueryMapWithTx(r.Context(), tx, h.Transform, s.Query, s.Params...)
		if er1 != nil {
			handleError(w, r, 500, er1.Error(), h.Error, er1)
			return
		}
		commit := ps.Get("commit")
		if commit == "true" {
			er3 := tx.Commit()
			if er3 != nil {
				handleError(w, r, http.StatusInternalServerError, er3.Error(), h.Error, er3)
				return
			}
			h.Cache.Remove(stx)
		}
		if len(res) > 0 {
			respond(w, http.StatusOK, res[0])
			return
		}
		respond(w, http.StatusOK, nil)
	}
}
func (h *Handler) ExecBatch(w http.ResponseWriter, r *http.Request) {
	var s []q.JStatement
	b := make([]q.Statement, 0)
	er0 := json.NewDecoder(r.Body).Decode(&s)
	if er0 != nil {
		http.Error(w, er0.Error(), http.StatusBadRequest)
		return
	}
	l := len(s)
	for i := 0; i < l; i++ {
		st := q.Statement{}
		st.Query = s[i].Query
		st.Params = q.ParseDates(s[i].Params, s[i].Dates)
		b = append(b, st)
	}
	ps := r.URL.Query()
	stx := ps.Get("tx")
	var er1 error
	var res int64
	if len(stx) == 0 {
		master := ps.Get("master")
		if master == "true" {
			res, er1 = q.ExecuteBatch(r.Context(), h.DB, b, true, true)
		} else {
			res, er1 = q.ExecuteAll(r.Context(), h.DB, b...)
		}
	} else {
		tx, er0 := h.Cache.Get(stx)
		if er0 != nil {
			http.Error(w, er0.Error(), http.StatusInternalServerError)
			return
		}
		if tx == nil {
			http.Error(w, "cannot get tx from cache. Maybe tx got timeout", http.StatusInternalServerError)
			return
		}
		tc := false
		commit := ps.Get("commit")
		if commit == "true" {
			tc = true
		}
		res, er1 = q.ExecuteStatements(r.Context(), tx, tc, b...)
		if tc && er1 == nil {
			h.Cache.Remove(stx)
		}
	}
	if er1 != nil {
		handleError(w, r, 500, er1.Error(), h.Error, er1)
		return
	}
	respond(w, http.StatusOK, res)
}

func handleError(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string, ...map[string]interface{}), err error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	respond(w, code, result)
}
func respond(w http.ResponseWriter, code int, result interface{}) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if result == nil {
		w.Write([]byte("null"))
		return nil
	}
	err := json.NewEncoder(w).Encode(result)
	return err
}
