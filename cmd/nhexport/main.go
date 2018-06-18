package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mrd0ll4r/nhexport"
	"github.com/pkg/errors"
)

var (
	from     string
	to       string
	payments bool
	addr     string
)

func init() {
	yesterday := time.Now().UTC().Add(-24 * time.Hour).Format("2006-01-02")
	today := time.Now().UTC().Format("2006-01-02")

	flag.StringVar(&addr, "addr", "", "your bitcoin address (mandatory)")
	flag.StringVar(&from, "from", yesterday, "begin date (inclusive)")
	flag.StringVar(&to, "to", today, "end date (exclusive)")
	flag.BoolVar(&payments, "payments", false, "whether to export hashrates+history (default) or payments")
}

func main() {
	flag.Parse()

	if len(addr) == 0 {
		flag.Usage()
		fmt.Println("Need address.")
		os.Exit(2)
	}

	os.Exit(doStuff())
}

func doStuff() int {
	since, err := time.Parse("2006-01-02", from)
	if err != nil {
		fmt.Println(err)
		return 2
	}

	until, err := time.Parse("2006-01-02", to)
	if err != nil {
		fmt.Println(err)
		return 2
	}

	mode := "hashrate"
	if payments {
		mode = "payments"
	}

	fName := fmt.Sprintf("%s-%s-%s-%s.csv", from, to, addr, mode)
	f, err := os.Create(fName)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer func() {
		f.Close()
		if err != nil {
			os.Remove(fName)
		}
	}()

	if payments {
		err = getPayments(f, since, until)
	} else {
		err = getHistory(f, since, until)
	}

	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("Wrote to %s\n", fName)
	return 0
}

func getPayments(f io.Writer, since, until time.Time) error {
	fmt.Println("getting $$$payments$$$ ...")
	payments, err := nhexport.GetPayoutsSince(addr, since)
	if err != nil {
		return err
	}

	untilTs := until.Unix()

	w := csv.NewWriter(f)
	defer w.Flush()
	err = w.Write(nhexport.PaymentsCSVHeader())
	if err != nil {
		return errors.Wrap(err, "unable to write CSV")
	}

	for _, payment := range payments {
		if payment.Timestamp >= untilTs {
			continue
		}

		err = w.Write(payment.CSV())
		if err != nil {
			return errors.Wrap(err, "unable to write CSV")
		}
	}

	return nil
}

func getHistory(f io.Writer, since, until time.Time) error {
	fmt.Println("getting ~~hashrate + history~~ ...")
	history, err := nhexport.GetAlgorithmHistoriesSince(addr, since)
	if err != nil {
		return err
	}

	untilTs := until.Unix()

	w := csv.NewWriter(f)
	defer w.Flush()
	err = w.Write(nhexport.HistoryCSVHeader())
	if err != nil {
		return errors.Wrap(err, "unable to write CSV")
	}

	for _, algo := range history {
		for _, entry := range algo.Data {
			if entry.Timestamp >= untilTs {
				continue
			}

			data := []string{algo.Algorithm.String()}
			data = append(data, entry.CSV()...)

			err = w.Write(data)
			if err != nil {
				return errors.Wrap(err, "unable to write CSV")
			}
		}
	}

	return nil
}
