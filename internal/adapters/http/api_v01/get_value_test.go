package api_v01

import (
	"database/sql"
	"github.com/0db0/metric-server/internal/contracts"
	mock_usecases "github.com/0db0/metric-server/internal/mocks/contracts"
	mock_logger "github.com/0db0/metric-server/internal/mocks/pkg/logger"
	"github.com/0db0/metric-server/internal/models"
	"github.com/0db0/metric-server/internal/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricAdapterGetMetric(t *testing.T) {
	type mocks struct {
		c *mock_usecases.MockCollectUseCase
		g *mock_usecases.MockGiveUseCase
		l *mock_logger.MockInterface
	}
	type fields struct {
		c     contracts.CollectUseCase
		g     contracts.GiveUseCase
		log   logger.Interface
		mocks *mocks
	}
	type args struct {
		endpoint   string
		httpMethod string
		payload    string
	}
	type setupMock func(m *mocks)
	type want struct {
		httpCode     int
		responseBody string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gMock := mock_usecases.NewMockGiveUseCase(ctrl)
	lMock := mock_logger.NewMockInterface(ctrl)
	m := &mocks{
		g: gMock,
		l: lMock,
	}
	tests := []struct {
		name   string
		args   args
		setup  setupMock
		fields fields
		want   want
	}{
		{
			name: "GET value 200",
			args: args{
				endpoint:   "/v0.1/value/counter/PollCount",
				httpMethod: http.MethodGet,
			},
			setup: func(m *mocks) {
				m.g.EXPECT().
					GetValue(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.Metric{
						ID:   1,
						Name: "PollCount",
						Type: "counter",
						Delta: sql.NullInt64{
							Int64: 74,
							Valid: true,
						},
					}, nil)
				m.l.EXPECT().Error(gomock.Any()).Times(0)
			},
			fields: fields{
				g:     gMock,
				log:   lMock,
				mocks: m,
			},
			want: want{
				httpCode:     http.StatusOK,
				responseBody: `{"name":"PollCount","type":"counter","delta":74}`,
			},
		},
		{
			name: "GET value 400 #1",
			args: args{
				endpoint:   "/v0.1/value",
				httpMethod: http.MethodGet,
			},
			setup: func(m *mocks) {
				m.g.EXPECT().
					GetValue(gomock.Any(), gomock.Any()).
					Times(0)
				m.l.EXPECT().
					Error("Error while building metric dto", gomock.Any()).
					Times(1)
			},
			fields: fields{
				g:     gMock,
				log:   lMock,
				mocks: m,
			},
			want: want{
				httpCode:     http.StatusBadRequest,
				responseBody: "",
			},
		},
		{
			name: "GET value 400 #2",
			args: args{
				endpoint:   "/v0.1/value/counter/",
				httpMethod: http.MethodGet,
			},
			setup: func(m *mocks) {
				m.g.EXPECT().
					GetValue(gomock.Any(), gomock.Any()).
					Times(0)
				m.l.EXPECT().
					Error("Error while building metric dto", gomock.Any()).
					Times(1)
			},
			fields: fields{
				g:     gMock,
				log:   lMock,
				mocks: m,
			},
			want: want{
				httpCode:     http.StatusBadRequest,
				responseBody: "",
			},
		},
		{
			name: "GET value 400 #2",
			args: args{
				endpoint:   "/v0.1/value//PollCount",
				httpMethod: http.MethodGet,
			},
			setup: func(m *mocks) {
				m.g.EXPECT().
					GetValue(gomock.Any(), gomock.Any()).
					Times(0)
				m.l.EXPECT().
					Error("Error while building metric dto", gomock.Any()).
					Times(1)
			},
			fields: fields{
				g:     gMock,
				log:   lMock,
				mocks: m,
			},
			want: want{
				httpCode:     http.StatusBadRequest,
				responseBody: "",
			},
		},
		{
			name: "GET value 404",
			args: args{
				endpoint:   "/v0.1/value/counter/UnknownMetric",
				httpMethod: http.MethodGet,
			},
			setup: func(m *mocks) {
				m.g.EXPECT().
					GetValue(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.Metric{}, sql.ErrNoRows)
				m.l.EXPECT().
					Error(gomock.Any(), gomock.Any()).
					Times(0)
			},
			fields: fields{
				g:     gMock,
				log:   lMock,
				mocks: m,
			},
			want: want{
				httpCode:     http.StatusNotFound,
				responseBody: "",
			},
		},
		{
			name: "GET value 500",
			args: args{
				endpoint:   "/v0.1/value/counter/PollCount",
				httpMethod: http.MethodGet,
			},
			setup: func(m *mocks) {
				m.g.EXPECT().
					GetValue(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.Metric{}, errors.New("unexpected error"))
				m.l.EXPECT().
					Error(gomock.Any(), gomock.Any()).
					Times(1)
			},
			fields: fields{
				g:     gMock,
				log:   lMock,
				mocks: m,
			},
			want: want{
				httpCode:     http.StatusInternalServerError,
				responseBody: "",
			},
		},
		{
			name: "POST value 200",
			args: args{
				endpoint:   "/v0.1/value",
				httpMethod: http.MethodPost,
				payload:    `{"id":"PollCount","type":"counter"}`,
			},
			setup: func(m *mocks) {
				m.g.EXPECT().
					GetValue(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.Metric{
						ID:   1,
						Name: "PollCount",
						Type: "counter",
						Delta: sql.NullInt64{
							Int64: 74,
							Valid: true,
						},
					}, nil)
				m.l.EXPECT().Error(gomock.Any()).Times(0)
			},
			fields: fields{
				g:     gMock,
				log:   lMock,
				mocks: m,
			},
			want: want{
				httpCode:     http.StatusOK,
				responseBody: `{"name":"PollCount","type":"counter","delta":74}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.args.payload)
			r := httptest.NewRequest(tt.args.httpMethod, tt.args.endpoint, body)
			w := httptest.NewRecorder()

			tt.setup(tt.fields.mocks)

			a := MetricAdapter{
				g:   tt.fields.g,
				log: tt.fields.log,
			}

			a.GetMetric(w, r)

			result := w.Result()
			defer result.Body.Close()

			actual, _ := io.ReadAll(result.Body)

			require.Equal(t, tt.want.httpCode, result.StatusCode)
			require.Equal(t, tt.want.responseBody, strings.TrimRight(string(actual), "\n"))
		})
	}
}
