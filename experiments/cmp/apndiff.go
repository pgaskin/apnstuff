package main

import (
	"bytes"
	"cmp"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"html"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
)

// go run apndiff.go google.xml.mcc302 motorola.xml.mcc302 samsung.xml.mcc302 > mcc302.html

func main() {
	const splitAPNTypes = false // set this to true to split apn types into separate rows

	var (
		files    = os.Args[1:]
		fileAPNs = make([][]APN, len(files))
	)
	for i, name := range files {
		apns, err := readAPNs(name, splitAPNTypes)
		if err != nil {
			panic(fmt.Errorf("%s: %w", name, err))
		}
		fileAPNs[i] = apns
	}

	type Key struct {
		MCC           string
		MNC           string
		MVNOType      string
		MVNOMatchData string
		Type          string
	}
	groupedAPNs := map[Key][]APN{} // [Key][file]
	for i, apns := range fileAPNs {
		for _, apn := range apns {
			key := Key{
				MCC:           apn["mcc"],
				MNC:           apn["mnc"],
				MVNOType:      apn["mvno_type"],
				MVNOMatchData: apn["mvno_match_data"],
				Type:          apn["type"],
			}
			if _, ok := groupedAPNs[key]; !ok {
				groupedAPNs[key] = make([]APN, len(files))
			}
			groupedAPNs[key][i] = apn
		}
	}

	fmt.Println(`<!DOCTYPE html>`)
	fmt.Println(`<html lang="en">`)
	fmt.Println(`<head>`)
	fmt.Println(`  <meta charset="utf-8">`)
	fmt.Println(`  <title>APN Diff</title>`)
	fmt.Println(`  <style>`)
	fmt.Println(`    html { font-size: 16px; }`)
	fmt.Println(`    body { font-family: sans-serif; font-size: 0.875rem; }`)
	fmt.Println(`    table { border-collapse: collapse; border: 1px solid #ccc; }`)
	fmt.Println(`    thead { position: sticky; top: 0; background: #eee; }`)
	fmt.Println(`    tr { border-bottom: 1px solid #ccc; }`)
	fmt.Println(`    th, td { padding: .25rem .5em; }`)
	fmt.Println(`    th { vertical-align: middle; text-align: left; }`)
	fmt.Println(`    td { vertical-align: top; }`)
	fmt.Println(`    td:nth-child(2) { border-right: 1px solid #ccc; }`)
	fmt.Println(`  </style>`)
	fmt.Println(`</head>`)
	fmt.Println(`<body>`)
	fmt.Println(`  <table>`)
	fmt.Println(`    <thead>`)
	fmt.Println(`      <tr>`)
	fmt.Println(`        <th>Match</th>`)
	fmt.Println(`        <th>Type</th>`)
	for _, name := range files {
		fmt.Printf("        <th>%s</th>\n", html.EscapeString(name))
	}
	fmt.Println(`      </tr>`)
	fmt.Println(`    </thead>`)
	fmt.Println(`    <tbody>`)
	for _, key := range slices.SortedStableFunc(maps.Keys(groupedAPNs), func(k1, k2 Key) int {
		return cmp.Or(
			cmp.Compare(k1.MCC, k2.MCC),
			cmp.Compare(k1.MNC, k2.MNC),
			cmp.Compare(k1.MVNOType, k2.MVNOType),
			cmp.Compare(k1.MVNOMatchData, k2.MVNOMatchData),
		)
	}) {
		fmt.Printf("      <tr style=\"background-color:%s\">\n", stringToBackgroundColor(key.MCC+key.MNC+key.MVNOType+key.MVNOMatchData))
		fmt.Printf("        <td>")
		fmt.Printf("%s %s", html.EscapeString(key.MCC), html.EscapeString(key.MNC))
		if key.MVNOType != "" {
			fmt.Printf(" %s=%s", html.EscapeString(key.MVNOType), html.EscapeString(key.MVNOMatchData))
		}
		fmt.Printf("</td>\n")
		fmt.Printf("        <td>%s</td>\n", html.EscapeString(key.Type))
		for _, apn := range groupedAPNs[key] {
			fmt.Printf("        <td>")
			fmt.Printf("<b>%s</b><br>", html.EscapeString(apn["apn"]))
			for _, attr := range apnAttrs {
				if val, ok := apn[attr]; ok {
					switch attr {
					case "mcc", "mnc", "mvno_type", "mvno_match_data", "type", "apn":
						continue // already grouped by these
					}
					fmt.Printf("%s = %s<br>", html.EscapeString(attr), html.EscapeString(val))
				}
			}
			fmt.Printf("</td>\n")
		}
		fmt.Printf("      </tr>\n")
	}
	fmt.Println(`    </tbody>`)
	fmt.Println(`  </table>`)
	fmt.Println(`</body>`)
}

func stringToBackgroundColor(str string) string {
	u := crc32.ChecksumIEEE([]byte(str))
	m := byte(0b1110_0000)
	return fmt.Sprintf("#%02x%02x%02x", byte(u)|m, byte(u>>8)|m, byte(u>>16)|m)
}

type APN map[string]string

var apnAttrs = []string{
	"carrier",
	"mcc",
	"mnc",
	"apn",
	"proxy",
	"port",
	"mmsc",
	"mmsproxy",
	"mmsport",
	"user",
	"password",
	"server",
	"authtype",
	"type",
	"protocol",
	"roaming_protocol",
	"carrier_enabled",
	"bearer_bitmask",
	"profile_id",
	"modem_cognitive",
	"max_conns",
	"wait_time",
	"max_conns_time",
	"mtu",
	"mvno_type",
	"mvno_match_data",
	"apn_set_id",
	"carrier_id",
	"skip_464xlat",
	"user_visible",
	"user_editable",
}

func readAPNs(name string, splitAPNTypes bool) ([]APN, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	// so we can run this on grepped configs
	if !bytes.HasSuffix(bytes.TrimSpace(buf), []byte("</apns>")) {
		buf = slices.Concat([]byte(`<apns>`), buf, []byte(`</apns>`))
	}

	// this doesn't do any validation other than that the xml is well-formed
	var (
		apns []APN
		dec  = xml.NewDecoder(bytes.NewReader(buf))
	)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			line, _ := dec.InputPos()
			if tok.Name == (xml.Name{Local: "apn"}) {
				apn := make(APN)
				for _, attr := range tok.Attr {
					// yes, this is inefficient
					if !slices.Contains(apnAttrs, attr.Name.Local) {
						return nil, fmt.Errorf("line %d: invalid apn attr %q (%q)", line, attr.Name.Local, tok)
					}
					apn[attr.Name.Local] = attr.Value
				}
				for _, a := range []string{"type", "mcc", "mnc"} {
					if _, ok := apn[a]; !ok {
						if a == "type" {
							apn[a] = "???wtf???" // some motorola ones are missing the type for some reason... is this valid, and if so, what's the default?
							continue
						}
						return nil, fmt.Errorf("line %d: missing apn attr %q (%q)", line, a, tok)
					}
				}
				if splitAPNTypes {
					for _, typ := range strings.Split(apn["type"], ",") {
						tmp := maps.Clone(apn)
						tmp["type"] = typ
						apns = append(apns, tmp)
					}
					continue
				}
				apn["type"] = strings.Join(slices.Sorted(slices.Values(strings.Split(apn["type"], ","))), ",") // so the grouping is consistent
				apns = append(apns, apn)
			}
		}
	}
	return apns, nil
}
