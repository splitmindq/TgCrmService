package bitrix

import (
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

func SendLeadToBitrix(log *slog.Logger, lead entities.LeadBitrix) error {
	urlWH := os.Getenv("BITRIX24_WEBHOOK")
	if urlWH == "" {
		return errors.New("BITRIX24_WEBHOOK is not set")
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
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("Bitrix response error", "status", resp.Status, "body", string(body))
		return fmt.Errorf("bitrix error: %s", resp.Status)
	}

	return nil
}
