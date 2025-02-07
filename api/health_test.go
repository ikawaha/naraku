package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-spectest/naraku/version"
	"github.com/go-spectest/spectest"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	"github.com/zeebo/assert"
)

var _ = Describe("Health api test", func() {
	var (
		t   GinkgoTInterface
		api *API
	)

	BeforeEach(func() {
		t = GinkgoT()
		api = NewAPI()
	})

	Context("Success to get health data", func() {
		BeforeSuite(func() {
			version.Version = "v0.0.1"
			version.Revision = "revision-0.0.1"
		})
		defer func() {
			version.Version = ""
			version.Revision = ""
		}()

		It("get server name, version, revision", func() {
			spectest.New().
				Report(spectest.SequenceDiagram(documentDirPath())).
				CustomReportName("health_success").
				Handler(api).
				Get("/v1/health").
				Expect(t).
				Body(`{"name": "naraku", "revision": "revision-0.0.1", "version":"v0.0.1"}`).
				Status(http.StatusOK).
				End()
		})
	})
})

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	t.Run("return health response", func(t *testing.T) {
		t.Parallel()

		e := echo.New()
		ctrl := NewHealthController()
		e.GET("/health", ctrl.health)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := ctrl.health(c) // Kick handler

		// Check response
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Check json response
		want := ctrl.newHealthResponse()
		expected, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())
	})
}

// This test is not parallel because it changes the global variable.
func TestNewHealthResponse(t *testing.T) { //nolint:paralleltest
	type args struct {
		tagVersion string
		revision   string
	}
	tests := []struct {
		name string
		args args
		want *HealthResponse
	}{
		{
			name: "unknown version",
			args: args{
				tagVersion: "",
				revision:   "rev",
			},
			want: &HealthResponse{
				Name:     "naraku",
				Version:  "unknown",
				Revision: "rev",
			},
		},
		{
			name: "unknown revision",
			args: args{
				tagVersion: "v0.0.1",
				revision:   "",
			},
			want: &HealthResponse{
				Name:     "naraku",
				Version:  "v0.0.1",
				Revision: "unknown",
			},
		},
		{
			name: "unknown version and revision",
			args: args{
				tagVersion: "",
				revision:   "",
			},
			want: &HealthResponse{
				Name:     "naraku",
				Version:  "unknown",
				Revision: "unknown",
			},
		},
		{
			name: "Get version and revision",
			args: args{
				tagVersion: "v0.0.0",
				revision:   "rev",
			},
			want: &HealthResponse{
				Name:     "naraku",
				Version:  "v0.0.0",
				Revision: "rev",
			},
		},
	}
	for _, tt := range tests { //nolint:paralleltest
		version.Version = tt.args.tagVersion
		version.Revision = tt.args.revision

		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewHealthController()
			if got := ctrl.newHealthResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHealthResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
