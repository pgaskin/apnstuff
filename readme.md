This repository contains stuff related to APNs on Android, especially where it relates to Canadian carriers.

The original goal of all this was fixing broken APNs for Freedom Mobile on LineageOS, and that's when I started learning about most of this.

- [mcc302.xml](./mcc302.xml) - APNs for Canadian carriers hand-selected from various sources.
- [aosp_notes.md](./aosp_notes.md) - Notes and links to source code for APN and CarrierConfig stuff (WIP).
- [apndiff.go](./experiments/cmp/apndiff.go) - Very rough tool to make comparing APN XMLs easier (note that this doesn't currently handle the fallback and carrier_id matching exactly the same way as AOSP; it's a rough approximation).
- [carrierId.proto](./android14-qpr3-release/carrierId.proto) - Protobuf schema for carrier_list.pb (shipped with AOSP TelephonyProvider and Google Carrier Services).
- [carrier_settings.proto](./android14-qpr3-release/carrier_settings.proto) - Protobuf schema for Google Carrier Services / Pixel carrierconfig.
