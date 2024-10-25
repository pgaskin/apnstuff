package main

import (
	"bytes"
	"cmp"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/crc32"
	"html"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/pgaskin/apnstuff/carriersettings/experiments/cmp/android14-qpr3-release/carrierid"
	"google.golang.org/protobuf/proto"
)

// curl https://android.googlesource.com/platform/packages/providers/TelephonyProvider/+/master/assets/latest_carrier_id/carrier_list.pb?format=TEXT | base64 -d > carrier_list.pb

// go run . 0 aosp.xml lineage.xml google.xml motorola.xml samsung.xml > mcc302-all.html
// go run . 0 google.xml motorola.xml samsung.xml > mcc302-good.html
// go run . 1 aosp.xml lineage.xml google.xml motorola.xml samsung.xml > mcc302-all.split.html
// go run . 1 google.xml motorola.xml samsung.xml > mcc302-good.split.html

func main() {
	splitAPNTypes, err := strconv.ParseBool(os.Args[1])
	if err != nil {
		panic(err)
	}

	var (
		files    = os.Args[2:]
		fileAPNs = make([][]APN, len(files))
	)
	for i, name := range files {
		apns, err := readAPNs(name, splitAPNTypes, "302")
		if err != nil {
			panic(fmt.Errorf("%s: %w", name, err))
		}
		fileAPNs[i] = apns
	}

	var carrierList carrierid.CarrierList
	{
		if buf, err := os.ReadFile("carrier_list.pb"); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				panic(err)
			}
		} else if err := proto.Unmarshal(buf, &carrierList); err != nil {
			panic(err)
		}
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
	fmt.Println(`    body { font-family: sans-serif; font-size: 0.75rem; }`)
	fmt.Println(`    table { border-collapse: collapse; border: 1px solid #ccc; }`)
	fmt.Println(`    thead { position: sticky; top: 0; background: #eee; }`)
	fmt.Println(`    tr { border-bottom: 1px solid #ccc; }`)
	fmt.Println(`    th, td { padding: .25rem .5em; }`)
	fmt.Println(`    th { vertical-align: middle; text-align: left; }`)
	fmt.Println(`    td { vertical-align: top; }`)
	fmt.Println(`    td:nth-child(1) { white-space: nowrap; }`)
	fmt.Println(`    td:nth-child(3) { border-right: 1px solid #ccc; }`)
	fmt.Println(`  </style>`)
	fmt.Println(`</head>`)
	fmt.Println(`<body>`)
	fmt.Println(`  <table>`)
	fmt.Println(`    <thead>`)
	fmt.Println(`      <tr>`)
	fmt.Println(`        <th>Match</th>`)
	fmt.Println(`        <th>Type</th>`)
	for _, name := range files {
		fmt.Printf("        <th>%s</th>\n", html.EscapeString(strings.TrimSuffix(name, ".xml")))
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
			cmp.Compare(k1.Type, k2.Type),
		)
	}) {
		fmt.Printf("      <tr style=\"background-color:%s\">\n", stringToBackgroundColor(key.MCC+key.MNC+key.MVNOType+key.MVNOMatchData))
		fmt.Printf("        <td>")
		fmt.Printf("<b>%s %s", html.EscapeString(key.MCC), html.EscapeString(key.MNC))
		if key.MVNOType != "" {
			fmt.Printf(" %s=%s", html.EscapeString(key.MVNOType), html.EscapeString(key.MVNOMatchData))
		}
		fmt.Printf("</b>")
		if carrierList.CarrierId != nil {
			if i := slices.IndexFunc(carrierList.CarrierId, func(e *carrierid.CarrierId) bool {
				return slices.ContainsFunc(e.CarrierAttribute, func(attr *carrierid.CarrierAttribute) bool {
					if !slices.Contains(attr.MccmncTuple, key.MCC+key.MNC) {
						return false
					}
					var hasMVNOCond bool
					if len(attr.ImsiPrefixXpattern) != 0 {
						hasMVNOCond = true
						if !strings.EqualFold(key.MVNOType, "imsi") || !containsStringFold(attr.ImsiPrefixXpattern, key.MVNOMatchData) {
							return false
						}
					}
					if len(attr.Spn) != 0 {
						hasMVNOCond = true
						if !strings.EqualFold(key.MVNOType, "spn") || !containsStringFold(attr.Spn, key.MVNOMatchData) {
							return false
						}
					}
					if len(attr.Plmn) != 0 {
						hasMVNOCond = true
						return false // no xml representation
					}
					if len(attr.Gid1) != 0 {
						hasMVNOCond = true
						if !strings.EqualFold(key.MVNOType, "gid") || !containsStringFold(attr.Gid1, key.MVNOMatchData) {
							return false
						}
					}
					if len(attr.Gid2) != 0 {
						hasMVNOCond = true
						return false // no xml representation
					}
					if len(attr.PreferredApn) != 0 {
						hasMVNOCond = true
						return false // no xml representation
					}
					if len(attr.IccidPrefix) != 0 {
						hasMVNOCond = true
						if !strings.EqualFold(key.MVNOType, "iccid") || !containsStringFold(attr.IccidPrefix, key.MVNOMatchData) {
							return false
						}
					}
					if key.MVNOType != "" && !hasMVNOCond {
						return false
					}
					return true
				})
			}); i != -1 {
				cid := carrierList.CarrierId[i]
				fmt.Printf("<br><br><span style=\"font-size:.75em\">full carrier_list match<br>")
				if cid.CarrierName != nil {
					fmt.Printf("<i>%s</i>", *cid.CarrierName)
				}
				if cid.CanonicalId != nil {
					fmt.Printf("<i> (%d)</i>", *cid.CanonicalId)
				}
				fmt.Printf("</span>")
			} else if i := slices.IndexFunc(carrierList.CarrierId, func(e *carrierid.CarrierId) bool {
				return slices.ContainsFunc(e.CarrierAttribute, func(attr *carrierid.CarrierAttribute) bool {
					if slices.Contains(attr.MccmncTuple, key.MCC+key.MNC) {
						switch key.MVNOType {
						case "imsi":
							return containsStringFold(attr.ImsiPrefixXpattern, key.MVNOMatchData)
						case "spn":
							return containsStringFold(attr.Spn, key.MVNOMatchData)
						case "gid":
							return containsStringFold(attr.Gid1, key.MVNOMatchData)
						case "iccid":
							return containsStringFold(attr.IccidPrefix, key.MVNOMatchData)
						}
					}
					return false
				})
			}); i != -1 {
				cid := carrierList.CarrierId[i]
				fmt.Printf("<br><br><span style=\"font-size:.75em\">partial carrier_list match<br>") // because the apn xml isn't capable of having multiple mvno match conditions
				if cid.CarrierName != nil {
					fmt.Printf("<i>%s</i>", *cid.CarrierName)
				}
				if cid.CanonicalId != nil {
					fmt.Printf("<i> (%d)</i>", *cid.CanonicalId)
				}
				fmt.Printf("</span>")
			}
		}
		fmt.Printf("</td>\n")
		fmt.Printf("        <td><b>%s</b></td>\n", html.EscapeString(key.Type))
		for i, apn := range groupedAPNs[key] {
			var isFallback bool
			if apn == nil {
				// fall back to the non-mvno-specific match if no apn for the mvno
				isFallback = true
				key.MVNOMatchData = ""
				key.MVNOType = ""
				if m, ok := groupedAPNs[key]; ok {
					apn = m[i]
				}
			}
			if isFallback {
				fmt.Printf("        <td style=\"opacity: 0.5\">")
			} else {
				fmt.Printf("        <td>")
			}
			if apn != nil {
				fmt.Printf("<b>apn = %s</b><br>", html.EscapeString(apn["apn"]))
				for _, attr := range apnAttrs {
					if val, ok := apn[attr]; ok {
						switch attr {
						case "mcc", "mnc", "mvno_type", "mvno_match_data", "type", "apn":
							continue // already grouped by these
						}
						fmt.Printf("%s = %s<br>", html.EscapeString(attr), html.EscapeString(val))
					}
				}
				if isFallback {
					fmt.Printf("<br><i>(fallback since no mvno-specific apn)</i>")
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

func containsStringFold(a []string, b string) bool {
	return slices.ContainsFunc(a, func(x string) bool {
		return strings.EqualFold(x, b)
	})
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

func readAPNs(name string, splitAPNTypes bool, filterMCC string) ([]APN, error) {
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
				if filterMCC != "" && !slices.ContainsFunc(tok.Attr, func(e xml.Attr) bool { return e.Name.Local == "mcc" && e.Value == filterMCC }) {
					continue
				}
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
