package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/libdns/digitalocean"
	"github.com/libdns/libdns"
)

func main() {
	token := os.Getenv("DO_AUTH_TOKEN")
	if token == "" {
		fmt.Printf("DO_AUTH_TOKEN not set\n")
		return
	}
	zone := os.Getenv("ZONE")
	if zone == "" {
		fmt.Printf("ZONE not set\n")
		return
	}
	provider := digitalocean.Provider{APIToken: token}

	records, err := provider.GetRecords(context.TODO(), zone)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	err = TXT_Test(provider, zone, records)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	err = A_Test(provider, zone, records)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func TXT_Test(provider digitalocean.Provider, zone string, records []libdns.Record) error {
	testName := "libdns-txt-test"
	testId := ""
	for _, record := range records {
		fmt.Printf("%s (.%s): %s, %s\n", record.Name, zone, record.Value, record.Type)
		if record.Name == testName {
			testId = record.ID
		}

	}

	if testId != "" {
		fmt.Printf("Replacing entry for %s\n", testName)
		_, err := provider.SetRecords(context.TODO(), zone, []libdns.Record{libdns.Record{
			Type:  "TXT",
			Name:  testName,
			Value: fmt.Sprintf("Replacement test entry created by libdns %s", time.Now()),
			TTL:   time.Duration(30) * time.Second,
			ID:    testId,
		}})
		return err
	}

	fmt.Printf("Creating new entry for %s\n", testName)
	_, err := provider.AppendRecords(context.TODO(), zone, []libdns.Record{libdns.Record{
		Type:  "TXT",
		Name:  testName,
		Value: fmt.Sprintf("This is a test entry created by libdns %s", time.Now()),
		TTL:   time.Duration(30) * time.Second,
	}})
	return err
}

func A_Test(provider digitalocean.Provider, zone string, records []libdns.Record) error {
	testName := "libdns-a-test"

	fmt.Printf("Creating new entry for %s\n", testName)
	_, err := provider.SetRecords(context.TODO(), zone, []libdns.Record{libdns.Record{
		Type:  "A",
		Name:  testName,
		Value: "127.0.0.1",
		TTL:   time.Duration(30) * time.Second,
	}})
	return err
}
