package bitrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lead-bitrix/entities"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type leadId struct {
	Result int `json:"result"`
}

func SendLeadToBitrix(log *slog.Logger, lead entities.JsonForm) (int, error) {
	urlWH := os.Getenv("BITRIX24ADD_WEBHOOK")
	if urlWH == "" {
		return 0, errors.New("BITRIX24ADD_WEBHOOK is not set")
	}

	data := url.Values{}
	data.Set("FIELDS[TITLE]", "Новый лид")
	data.Set("FIELDS[NAME]", lead.Name)
	data.Set("FIELDS[LAST_NAME]", "-")
	data.Set("FIELDS[EMAIL][0][VALUE]", lead.Email)
	data.Set("FIELDS[EMAIL][0][VALUE_TYPE]", "WORK")
	data.Set("FIELDS[PHONE][0][VALUE]", lead.Phone)
	data.Set("FIELDS[PHONE][0][VALUE_TYPE]", "WORK")
	data.Set("FIELDS[SOURCE_DESCRIPTION]", lead.Source)

	resp, err := http.Post(
		urlWH,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		log.Error("SendLeadToBitrix POST error:", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Error("Bitrix response error", "status", resp.Status, "body", string(body))
		return 0, fmt.Errorf("bitrix error: %s", resp.Status)
	}

	var bitrixId leadId

	if err := json.Unmarshal(body, &bitrixId); err != nil {
		log.Error("Bitrix response error", "error", err, "body", string(body))
		return 0, err
	}
	return bitrixId.Result, nil
}
