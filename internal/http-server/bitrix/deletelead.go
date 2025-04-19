package bitrix

import (
	"errors"
	"lead-bitrix/entities"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func DeleteLead(log *slog.Logger, lead *entities.Lead) error {

	urlWH := os.Getenv("BITRIX24DELETE_WEBHOOK")
	if urlWH == "" {
		return errors.New("BITRIX24ADD_WEBHOOK is not set")
	}

	data := url.Values{}
	data.Set("id", strconv.Itoa(lead.Id))

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

	return nil
}
