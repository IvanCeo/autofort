package wordorder

import (
	"autofort/internal/usecase"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type WorkOrderClient struct {
	baseURL string
	http    *http.Client
}

func NewWorkOrderClient(ctx context.Context, timeout time.Duration) (*WorkOrderClient, error) {
	env := os.Getenv("APP_ENV")
	isProd := env == "PROD"

	hostKey := "WORK_ORDER_HOST_DEV"
	portKey := "WORK_ORDER_PORT_DEV"
	if isProd {
		hostKey = "WORK_ORDER_HOST_PROD"
		portKey = "WORK_ORDER_PORT_PROD"
	}

	host := os.Getenv(hostKey)
	port := os.Getenv(portKey)

	if host == "" {
		return nil, fmt.Errorf("%s is empty", hostKey)
	}
	if port == "" {
		return nil, fmt.Errorf("%s is empty", portKey)
	}

	baseURL := fmt.Sprintf("http://%s:%s", host, port)

	return &WorkOrderClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *WorkOrderClient) CreateWorkOrderPDF(req *usecase.WorkOrderCreateRequest) (*usecase.WorkOrderCreateResponse, error) {
	// Лучше принимать ctx от usecase/handler, но пока держим просто Background.
	// Если захочешь — я покажу, как протянуть ctx через слои.
	ctx, cancel := context.WithTimeout(context.Background(), c.http.Timeout)
	defer cancel()

	// 1) JSON encode
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return &usecase.WorkOrderCreateResponse{}, err
	}

	// 2) HTTP request
	url := c.baseURL + "/work-order/pdf"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return &usecase.WorkOrderCreateResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/pdf")

	// 3) Do
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return &usecase.WorkOrderCreateResponse{}, err
	}
	defer resp.Body.Close()

	// 4) status handling
	if resp.StatusCode != http.StatusOK {
		// читаем кусочек тела, чтобы видеть ошибку сервиса (и логировать)
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if len(b) == 0 {
			return &usecase.WorkOrderCreateResponse{}, fmt.Errorf("work-order service returned %d", resp.StatusCode)
		}
		return &usecase.WorkOrderCreateResponse{}, fmt.Errorf("work-order service returned %d: %s", resp.StatusCode, string(b))
	}

	// 5) content-type check (мягкий)
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		// бывает, сервис не ставит — не валим
	} else if ct != "application/pdf" && ct != "application/pdf; charset=binary" {
		// тоже не всегда строго совпадает, но пусть будет защита от JSON-ответа
		// читаем чуть-чуть и возвращаем ошибку
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return &usecase.WorkOrderCreateResponse{}, errors.New("unexpected content-type from work-order service: " + ct + " body: " + string(b))
	}

	// 6) read pdf bytes
	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return &usecase.WorkOrderCreateResponse{}, err
	}
	if len(pdf) == 0 {
		return &usecase.WorkOrderCreateResponse{}, errors.New("empty pdf from work-order service")
	}

	return &usecase.WorkOrderCreateResponse{PDF: pdf}, nil
}
