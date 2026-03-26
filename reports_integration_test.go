//go:build integration

package onspring_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestReportsIntegration(t *testing.T) {
	loadEnvFile(t)
	client := createClient(t)
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		t.Run("should return a report", func(t *testing.T) {
			reportId := requireEnvInt(t, "TEST_REPORT")

			report, err := client.Reports.Get(ctx, reportId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if report.Columns == nil {
				t.Error("Expected columns to not be nil")
			}

			if report.Rows == nil {
				t.Error("Expected rows to not be nil")
			}

			for _, row := range report.Rows {
				if row.Cells == nil {
					t.Error("Expected cells to not be nil")
				}
			}
		})

		t.Run("should return report data for a report with chart data when report data is requested", func(t *testing.T) {
			reportId := requireEnvInt(t, "TEST_REPORT_WITH_CHART_DATA")

			report, err := client.Reports.Get(ctx, reportId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if report.Columns == nil {
				t.Error("Expected columns to not be nil")
			}

			if report.Rows == nil {
				t.Error("Expected rows to not be nil")
			}
		})

		t.Run("should return chart data for a report with a chart when chart data is requested", func(t *testing.T) {
			reportId := requireEnvInt(t, "TEST_REPORT_WITH_CHART_DATA")

			report, err := client.Reports.Get(
				ctx,
				reportId,
				onspring.WithDataFormat("Raw"),
				onspring.WithDataType("ChartData"),
			)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if report.Columns == nil {
				t.Error("Expected columns to not be nil")
			}

			if report.Rows == nil {
				t.Error("Expected rows to not be nil")
			}
		})

		t.Run("should return a 400 error if chart data is requested for a report without chart data", func(t *testing.T) {
			reportId := requireEnvInt(t, "TEST_REPORT")

			_, err := client.Reports.Get(
				ctx,
				reportId,
				onspring.WithDataFormat("Raw"),
				onspring.WithDataType("ChartData"),
			)

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 error if the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			reportId := requireEnvInt(t, "TEST_REPORT")

			_, err := invalidClient.Reports.Get(ctx, reportId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error if the api key does not have access to the report", func(t *testing.T) {
			reportIdNoAccess := requireEnvInt(t, "TEST_REPORT_NO_ACCESS")

			_, err := client.Reports.Get(ctx, reportIdNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 error if the report does not exist", func(t *testing.T) {
			_, err := client.Reports.Get(ctx, 0)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("should return a list of reports", func(t *testing.T) {
			appId := requireEnvInt(t, "TEST_APP_ID")

			page, err := client.Reports.List(ctx, appId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.PageNumber == 0 {
				t.Error("Expected page number to not be zero")
			}

			if page.TotalRecords == 0 {
				t.Error("Expected total records to not be zero")
			}

			if len(page.Items) == 0 {
				t.Error("Expected items to not be empty")
			}

			for _, report := range page.Items {
				if report.Id == 0 {
					t.Error("Expected report id to not be zero")
				}

				if report.AppId == 0 {
					t.Error("Expected report appId to not be zero")
				}

				if report.Name == "" {
					t.Error("Expected report name to not be empty")
				}
			}
		})

		t.Run("should return a 400 error if page size is invalid", func(t *testing.T) {
			appId := requireEnvInt(t, "TEST_APP_ID")

			_, err := client.Reports.List(ctx, appId, onspring.WithPageSize(1001))

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 error if the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			appId := requireEnvInt(t, "TEST_APP_ID")

			_, err := invalidClient.Reports.List(ctx, appId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error if the api key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")

			_, err := client.Reports.List(ctx, appIdNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("should iterate all reports for an app", func(t *testing.T) {
			appId := requireEnvInt(t, "TEST_APP_ID")

			var reports []onspring.Report

			for report, err := range client.Reports.ListAll(ctx, appId, onspring.WithPageSize(1)) {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				reports = append(reports, report)
			}

			if len(reports) == 0 {
				t.Error("Expected to iterate at least one report")
			}
		})
	})
}
