package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	helpUsage = "RangeRadar -range=\"10.0.0.0-10.255.255.255\" -output json"
)

func ipRangeToCIDRs(ipRangeStr string) ([]*net.IPNet, error) {
	ips := strings.Split(ipRangeStr, "-")
	if len(ips) != 2 {
		return nil, fmt.Errorf("invalid IP range format: %s", ipRangeStr)
	}

	startIP := net.ParseIP(ips[0])
	endIP := net.ParseIP(ips[1])
	if startIP == nil || endIP == nil {
		return nil, fmt.Errorf("invalid IP address in range: %s", ipRangeStr)
	}

	cidrs := []*net.IPNet{}
	for start := startIP; ; {
		mask := net.CIDRMask(32, 32)
		for ones := 0; ones <= 32; ones++ {
			cidr := net.IPNet{IP: start, Mask: mask}
			if !cidr.Contains(endIP) {
				break
			}
			mask = net.CIDRMask(32-ones, 32)
		}

		cidrs = append(cidrs, &net.IPNet{IP: start, Mask: mask})
		lastIP := make(net.IP, len(start))
		copy(lastIP, start)
		for i := len(lastIP) - 1; i >= 0; i-- {
			lastIP[i] |= ^mask[i]
			if lastIP[i] != 0xff {
				break
			}
		}

		if lastIP.Equal(endIP) {
			break
		}

		start = make(net.IP, len(lastIP))
		copy(start, lastIP)
		for i := len(start) - 1; i >= 0; i-- {
			start[i]++
			if start[i] != 0 {
				break
			}
		}
	}

	return cidrs, nil
}

func main() {
	var ipRangeStr string
	var outputFormat string
	flag.StringVar(&ipRangeStr, "range", "", "an IP address range to convert to CIDR blocks")
	flag.StringVar(&outputFormat, "output", "terminal", "the output format (json, csv, or terminal)")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Println("Convert an IP address range to a list of CIDR blocks")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println(helpUsage)
	}

	flag.Parse()

	if ipRangeStr == "" {
		fmt.Println("Error: -range is required.")
		flag.Usage()
		os.Exit(1)
	}

	startTime := time.Now()
	cidrs, err := ipRangeToCIDRs(ipRangeStr)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	switch outputFormat {
	case "json":
		filename := fmt.Sprintf("cidrs_%s_%s.json", strings.Replace(ipRangeStr, "/", "-", -1), time.Now().Format("2006-01-02T15-04-05"))
		err = outputJSON(cidrs, filename)
		if err != nil {
			fmt.Printf("Error writing JSON output to file: %v\n", err)
		}
	case "csv":
		filename := fmt.Sprintf("cidrs_%s_%s.csv", strings.Replace(ipRangeStr, "/", "-", -1), time.Now().Format("2006-01-02T15-04-05"))
		err := outputCSV(cidrs, filename)
		if err != nil {
			fmt.Printf("Error writing CSV file: %v\n", err)
		}
	default:
		outputTerminal(cidrs)
	}

	fmt.Printf("Took %v seconds to complete.\n", time.Since(startTime).Seconds())
}

func outputJSON(cidrs []*net.IPNet, filename string) error {
	type CIDR struct {
		Range string `json:"range"`
	}
	var data []CIDR
	for _, cidr := range cidrs {
		data = append(data, CIDR{cidr.String()})
	}
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func outputCSV(cidrs []*net.IPNet, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, cidr := range cidrs {
		writer.Write([]string{cidr.String()})
	}

	return nil
}

func outputTerminal(cidrs []*net.IPNet) {
	for _, cidr := range cidrs {
		fmt.Println(cidr.String())
	}
}
