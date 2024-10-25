#!/bin/bash
set -euxo pipefail

for x in carrier_list carrier_settings carrierId; do
for y in android-12.0.0_r4 android-13.0.0_r3 android14-qpr3-release; do
protoc -I $y $x.proto -o $y/$x.proto.pb
done
done

rm -rf out
mkdir out

mkdir out/2021
bsdtar xOf raven_SD1A.210817.015.A4/TelephonyProvider.apk assets/carrier_list.pb | protoscope -descriptor-set android-12.0.0_r4/carrierId.proto.pb -message-type carrierIdentification.CarrierList -print-field-names -print-enum-names > out/2021/carrierId.txt
protoscope -descriptor-set android-12.0.0_r4/carrier_list.proto.pb -message-type com.google.carrier.CarrierList -print-field-names -print-enum-names < raven_SD1A.210817.015.A4/CarrierSettings/carrier_list.pb > out/2021/carrier_list.txt
for x in freedommobile bell telus shaw rogers videotron; do
protoscope -descriptor-set android-12.0.0_r4/carrier_settings.proto.pb -message-type com.google.carrier.CarrierSettings -print-field-names -print-enum-names < raven_SD1A.210817.015.A4/CarrierSettings/${x}_ca.pb > out/2021/carrier_settings_$x.txt
done
python3 lineage/scripts/carriersettings-extractor/carriersettings_extractor.py -i raven_SD1A.210817.015.A4/CarrierSettings -a out/2021 -v out/2021
xmllint --format out/2021/apns-conf.xml --output out/2021/apns-conf.xml

mkdir out/2022
bsdtar xOf panther_TD1A.220804.009.A2/TelephonyProvider.apk assets/carrier_list.pb | protoscope -descriptor-set android-13.0.0_r3/carrierId.proto.pb -message-type carrierIdentification.CarrierList -print-field-names -print-enum-names > out/2022/carrierId.txt
protoscope -descriptor-set android-13.0.0_r3/carrier_list.proto.pb -message-type com.google.carrier.CarrierList -print-field-names -print-enum-names < panther_TD1A.220804.009.A2/CarrierSettings/carrier_list.pb > out/2022/carrier_list.txt
for x in freedommobile bell telus shaw rogers videotron; do
protoscope -descriptor-set android-13.0.0_r3/carrier_settings.proto.pb -message-type com.google.carrier.CarrierSettings -print-field-names -print-enum-names < panther_TD1A.220804.009.A2/CarrierSettings/${x}_ca.pb > out/2022/carrier_settings_$x.txt
done
python3 lineage/scripts/carriersettings-extractor/carriersettings_extractor.py -i panther_TD1A.220804.009.A2/CarrierSettings -a out/2022 -v out/2022
xmllint --format out/2022/apns-conf.xml --output out/2022/apns-conf.xml

mkdir out/2024
bsdtar xOf shiba_AP2A.240905.003/TelephonyProvider.apk assets/carrier_list.pb | protoscope -descriptor-set android14-qpr3-release/carrierId.proto.pb -message-type carrierIdentification.CarrierList -print-field-names -print-enum-names > out/2024/carrierId.txt
protoscope -descriptor-set android14-qpr3-release/carrier_list.proto.pb -message-type com.google.carrier.CarrierList -print-field-names -print-enum-names < shiba_AP2A.240905.003/CarrierSettings/carrier_list.pb > out/2024/carrier_list.txt
for x in freedommobile bell telus shaw rogers videotron; do
protoscope -descriptor-set android14-qpr3-release/carrier_settings.proto.pb -message-type com.google.carrier.CarrierSettings -print-field-names -print-enum-names < shiba_AP2A.240905.003/CarrierSettings/${x}_ca.pb > out/2024/carrier_settings_$x.txt
done
python3 lineage/scripts/carriersettings-extractor/carriersettings_extractor.py -i shiba_AP2A.240905.003/CarrierSettings -a out/2024 -v out/2024
xmllint --format out/2024/apns-conf.xml --output out/2024/apns-conf.xml

exit

### to install and apply new apns config without rebuilding lineage:
adb root

# put the new apns-conf.xml in the ota-updated apns location
# (which takes precedence over the /product/etc/apns-conf.xml which lineage ships)
adb push out/2024/apns-conf.xml /data/misc/apns/apns-conf.xml

# edit the saved build-id so the telephony provider regenerates the carriers table in telephony.db
adb shell "sed -i -E 's/(<string name=.ro_build_id.>)[^<]+(<.string>)/\1dummy\2/g' /data/user_de/0/com.android.providers.telephony/shared_prefs/build-id.xml"

# reboot
adb reboot

# reset network settings in UI
# also need to reset apn setting (via the three dot menu) specifically due to some lineage oddness, or old ones may get left behind

# reboot
adb reboot
