package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"reddit_challenge/model"
)

type MockDb struct {
	tech []*model.Technology
	err  error
}

func (m *MockDb) GetStats() ([]*model.Technology, error) {
	return m.tech, m.err
}

func TestApp_GetStats(t *testing.T) {
	app := App{
		d: &MockDb{
			tech: []*model.Technology{
				{"Tech1", "Details1"},
				{"Tech2", "Details2"},
			},
		},
	}

	r, _ := http.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()

	app.GetStats(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	want := `[{"name":"Tech1","details":"Details1"},{"name":"Tech2","details":"Details2"}]` + "\n"
	if got := w.Body.String(); got != want {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

func TestApp_GetStats_WithDBError(t *testing.T) {
	app := App{
		d: &MockDb{
			tech: nil,
			err:  errors.New("unknown error"),
		},
	}

	r, _ := http.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()

	app.GetStats(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}
