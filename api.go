package nhexport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// A Payment describes a singular payment.
type Payment struct {
	Amount    decimal.Decimal `json:"amount"`
	Fee       decimal.Decimal `json:"fee"`
	TxID      string          `json:"TXID"`
	Timestamp int64           `json:"time"`
	Type      int             `json:"type"`
}

// PaymentsCSVHeader returns the CSV header for the payments "table".
func PaymentsCSVHeader() []string {
	return []string{
		"timestamp",
		"amount",
		"fee",
		"transactionID",
		"type"}
}

// CSV returns the fields of the Payment for CSV printing.
func (e *Payment) CSV() []string {
	amount := e.Amount.String()
	fee := e.Fee.String()
	txID := e.TxID
	ts := time.Unix(e.Timestamp, 0).UTC().Format(time.RFC3339)
	var typ string
	switch e.Type {
	case 0:
		typ = "standard"
	case 1:
		typ = "internal"
	default:
		typ = "unknown"
	}

	return []string{
		ts,
		amount,
		fee,
		txID,
		typ}
}

// A HistoryEntry describes the hashrate and unpaid balance of one algorithm at
// one point in time.
type HistoryEntry struct {
	Timestamp     int64
	Hashrate      HashrateEntry
	UnpaidBalance decimal.Decimal
}

// HistoryCSVHeader returns the headers of the algorithm history "table".
func HistoryCSVHeader() []string {
	return []string{
		"algorithm",
		"timestamp",
		"accepted",
		"rejectedStale",
		"rejectedTarget",
		"rejectedDuplicate",
		"rejectedOther",
		"unpaidBalance"}
}

// CSV returns the fields of a HistoryEntry for CSV printing.
func (e *HistoryEntry) CSV() []string {
	ts := time.Unix(e.Timestamp, 0).UTC().Format(time.RFC3339)
	accepted := e.Hashrate.Accepted.String()
	stale := e.Hashrate.RejectedStale.String()
	target := e.Hashrate.RejectedTarget.String()
	duplicate := e.Hashrate.RejectedDuplicate.String()
	other := e.Hashrate.RejectedOther.String()
	unpaid := e.UnpaidBalance.String()

	return []string{
		ts,
		accepted,
		stale,
		target,
		duplicate,
		other,
		unpaid}
}

// A HashrateEntry describes hashrates.
type HashrateEntry struct {
	Accepted          decimal.Decimal `json:"a"`
	RejectedStale     decimal.Decimal `json:"rs"`
	RejectedTarget    decimal.Decimal `json:"rt"`
	RejectedDuplicate decimal.Decimal `json:"rd"`
	RejectedOther     decimal.Decimal `json:"ro"`
}

type rawAlgorithmHistory struct {
	Algorithm int                 `json:"algo"`
	Data      [][]json.RawMessage `json:"data"`
}

func (e *rawAlgorithmHistory) toAlgorithmHistory() (*AlgorithmHistory, error) {
	toReturn := AlgorithmHistory{
		Algorithm: Algorithm(e.Algorithm),
	}

	for _, entry := range e.Data {
		raw, err := json.Marshal(entry)
		if err != nil {
			panic(err)
		}
		if len(entry) != 3 {
			return nil, fmt.Errorf("invalid history entry, wanted 3 items, got %d, entry: %q", len(entry), string(raw))
		}

		raw0, _ := entry[0].MarshalJSON()
		raw1, _ := entry[1].MarshalJSON()
		raw2, _ := entry[2].MarshalJSON()

		newEntry := HistoryEntry{}

		err = json.Unmarshal(raw0, &newEntry.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode timestamp")
		}
		newEntry.Timestamp = newEntry.Timestamp * 300

		err = json.Unmarshal(raw1, &newEntry.Hashrate)
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode hashrate")
		}

		err = json.Unmarshal(raw2, &newEntry.UnpaidBalance)
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode unpaid balance")
		}

		toReturn.Data = append(toReturn.Data, newEntry)
	}

	return &toReturn, nil
}

// An AlgorithmHistory describes the history of one algorithm.
type AlgorithmHistory struct {
	Algorithm Algorithm
	Data      []HistoryEntry
}

type apiResponse struct {
	Method string                `json:"method"`
	Result providerStatsResponse `json:"result"`
}

type providerStatsResponse struct {
	Error    string                `json:"error"`
	Address  string                `json:"addr"`
	Payments []Payment             `json:"payments"`
	Past     []rawAlgorithmHistory `json:"past"`
}

func providerStatsEx(addr string, from time.Time) (*providerStatsResponse, error) {
	fromUnix := from.Unix()
	params := url.Values{}
	params.Set("method", "stats.provider.ex")
	params.Set("addr", addr)
	params.Set("from", fmt.Sprint(fromUnix))
	p := params.Encode()

	resp, err := http.Get("https://api.nicehash.com/api?" + p)
	if err != nil {
		return nil, errors.Wrap(err, "unable to perform request")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response body")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned status %d: %s, body %s", resp.StatusCode, resp.Status, string(b))
	}

	r := apiResponse{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode response (raw: %s)", string(b))
	}

	if len(r.Result.Error) != 0 {
		return nil, fmt.Errorf("API returned error: %q", r.Result.Error)
	}

	if r.Method != "stats.provider.ex" {
		return nil, fmt.Errorf("got result for wrong method? Expected %q, got %q, body %s", "stats.provider.ex", r.Method, string(b))
	}

	if r.Result.Address != addr {
		return nil, fmt.Errorf("got result for wrong address? Expected %q, got %q, body %s", addr, r.Result.Address, string(b))
	}

	return &r.Result, nil
}

// GetPayoutsSince returns payouts since the given date.
func GetPayoutsSince(addr string, from time.Time) ([]Payment, error) {
	resp, err := providerStatsEx(addr, from)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query API")
	}

	return resp.Payments, nil
}

// GetAlgorithmHistoriesSince returns all AlgorithmHistories since the given
// date.
func GetAlgorithmHistoriesSince(addr string, from time.Time) ([]AlgorithmHistory, error) {
	resp, err := providerStatsEx(addr, from)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query API")
	}

	var history []AlgorithmHistory
	for _, entry := range resp.Past {
		newEntry, err := entry.toAlgorithmHistory()
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode response")
		}

		history = append(history, *newEntry)
	}

	return history, nil
}
