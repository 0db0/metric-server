package api_v01

import (
	"github.com/0db0/metric-server/internal/contracts"
	"github.com/0db0/metric-server/internal/dto"
	"github.com/0db0/metric-server/internal/mocks/adapters/http/api_v01"
	mock_usecases "github.com/0db0/metric-server/internal/mocks/contracts"
	mock_logger "github.com/0db0/metric-server/internal/mocks/pkg/logger"
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

func TestMetricAdapter_Collect(t *testing.T) {
	type mocks struct {
		c  *mock_usecases.MockCollectUseCase
		rb *api_v01.MockRequestDtoBuilder
		l  *mock_logger.MockInterface
	}
	type fields struct {
		c   contracts.CollectUseCase
		g   contracts.GiveUseCase
		rb  RequestDtoBuilder
		log logger.Interface
	}
	type args struct {
		httpMethod string
		endpoint   string
		payload    string
	}
	type setupMocks func(m *mocks)
	type want struct {
		httpCode     int
		responseBody string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cMock := mock_usecases.NewMockCollectUseCase(ctrl)
	lMock := mock_logger.NewMockInterface(ctrl)
	rbMock := api_v01.NewMockRequestDtoBuilder(ctrl)
	m := &mocks{
		c:  cMock,
		l:  lMock,
		rb: rbMock,
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		setupMocks setupMocks
		mocks      *mocks
		want       want
	}{
		{
			name: "POST value via body 200",
			fields: fields{
				c:   cMock,
				rb:  rbMock,
				log: lMock,
			},
			mocks: m,
			setupMocks: func(m *mocks) {
				m.rb.EXPECT().
					CreateCollectDto(gomock.Any()).
					Times(1).
					Return(dto.CollectDto{}, nil)
				m.c.EXPECT().
					CollectOne(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
				m.l.EXPECT().Error(gomock.Any()).Times(0)
			},
			args: args{
				httpMethod: http.MethodPost,
				payload:    `{"id":"Alloc","type":"gauge","value":55.55}`,
				endpoint:   "/v0.1/update",
			},
			want: want{
				httpCode:     200,
				responseBody: "",
			},
		},
		{
			name: "POST value via path 200",
			fields: fields{
				c:   cMock,
				rb:  rbMock,
				log: lMock,
			},
			mocks: m,
			setupMocks: func(m *mocks) {
				m.rb.EXPECT().
					CreateCollectDto(gomock.Any()).
					Times(1).
					Return(dto.CollectDto{}, nil)
				m.c.EXPECT().
					CollectOne(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
				m.l.EXPECT().Error(gomock.Any()).Times(0)
			},
			args: args{
				httpMethod: http.MethodPost,
				payload:    "",
				endpoint:   "/v0.1/update/gauge/Alloc/55.5",
			},
			want: want{
				httpCode:     200,
				responseBody: "",
			},
		},
		{
			name: "POST value 400",
			fields: fields{
				c:   cMock,
				rb:  rbMock,
				log: lMock,
			},
			mocks: m,
			setupMocks: func(m *mocks) {
				m.rb.EXPECT().
					CreateCollectDto(gomock.Any()).
					Times(1).
					Return(dto.CollectDto{}, errors.New("error"))
				m.c.EXPECT().
					CollectOne(gomock.Any(), gomock.Any()).
					Times(0)
				m.l.EXPECT().Error(gomock.Any()).Times(0)
			},
			args: args{
				httpMethod: http.MethodPost,
				payload:    `{"id":"Alloc","type":"gauge","value":55.55}`,
				endpoint:   "/v0.1/update",
			},
			want: want{
				httpCode:     400,
				responseBody: "",
			},
		},
		{
			name: "POST value 500",
			fields: fields{
				c:   cMock,
				rb:  rbMock,
				log: lMock,
			},
			mocks: m,
			setupMocks: func(m *mocks) {
				m.rb.EXPECT().
					CreateCollectDto(gomock.Any()).
					Times(1).
					Return(dto.CollectDto{}, nil)
				m.c.EXPECT().
					CollectOne(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("error"))
				m.l.EXPECT().Error(gomock.Any()).Times(0)
			},
			args: args{
				httpMethod: http.MethodPost,
				payload:    `{"id":"Alloc","type":"gauge","value":55.55}`,
				endpoint:   "/v0.1/update",
			},
			want: want{
				httpCode:     500,
				responseBody: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.args.payload)
			r := httptest.NewRequest(tt.args.httpMethod, tt.args.endpoint, body)
			w := httptest.NewRecorder()

			a := MetricAdapter{
				c:   tt.fields.c,
				rb:  tt.fields.rb,
				log: tt.fields.log,
			}

			tt.setupMocks(tt.mocks)
			a.Collect(w, r)

			result := w.Result()
			defer result.Body.Close()

			actual, _ := io.ReadAll(result.Body)

			require.Equal(t, tt.want.httpCode, result.StatusCode)
			require.Equal(t, tt.want.responseBody, strings.TrimRight(string(actual), "\n"))
		})
	}
}
