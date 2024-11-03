# AOSP Notes

Some notes about how APNs actually work on Android since it's not obvious and documentation is nearly nonexistent. Mostly based on [`android14-qpr3-release`](https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:).

TODO: finish collecting my other notes

TODO: finish writing this all up

## CarrierID Matching

https://source.android.com/docs/core/connect/carrierid

TODO: sort order

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:frameworks/opt/telephony/src/java/com/android/internal/telephony/CarrierResolver.java;l=580-597

TODO: special case for gid matching

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:frameworks/opt/telephony/src/java/com/android/internal/telephony/CarrierResolver.java;l=752-758?q=apn%20match

TODO: specific vs non-specific carrier_id (i.e., mvno vs parent -- but note that many don't have this set and aosp usually just uses the specific one)

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:frameworks/opt/telephony/src/java/com/android/internal/telephony/CarrierResolver.java;l=1012-1052

TODO: carrier uicc app

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:frameworks/opt/telephony/src/java/com/android/internal/telephony/CarrierResolver.java;l=760-774

https://source.android.com/docs/core/connect/uicc

## APN Matching

TODO: telephony provider db

TODO: match order (any mccmnc+mvno, any mccmnc, + remaining carrier_id matches)

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=4499-4606

TODO: figure out what happens if type fields contains overlap

TODO: user preferred apns

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=3449-3551

TODO: user edited apns

TODO: effective apn list

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/apps/Settings/src/com/android/settings/network/apn/ApnSettings.java;l=322-394

## APN Settings

TODO: all fields

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2507-2554

TODO: bearer_bitmask values, lte_ca -> lte normalization

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=333-375

https://cs.android.com/android/platform/superproject/main/+/main:frameworks/base/telephony/java/android/telephony/TelephonyManager.java;l=14742-14912;drc=517003e1168f7d0e07cac9c60e67b73f0e28dbde;bpv=0;bpt=1

TODO: network_bitmask/bearer_bitmask sync

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=5624-5647

TODO: valid aosp type values, removal of spaces

https://cs.android.com/android/platform/superproject/main/+/main:packages/apps/Settings/src/com/android/settings/network/apn/ApnTypes.kt;l=32-48;drc=2faad6df6ee4a5bcd067c076326af19017f59def

TODO: valid mmsc, https?:// re

https://cs.android.com/android/platform/superproject/main/+/main:packages/apps/Settings/src/com/android/settings/network/apn/ApnStatus.kt;l=274-288;drc=57956a8bb8be4a5711af032911055562efed5c6c

TODO: misc validation

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2711-2827

## APN Load

TODO: load path order, ota receiver, telephony content provider, carrierconfig service

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:frameworks/base/services/core/java/com/android/server/updates/ApnDbInstallReceiver.java

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/services/Telephony/AndroidManifest.xml;l=584-591

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/modules/Permission/tests/cts/permissionpolicy/res/raw/android_manifest.xml;l=4605-4611

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=884-890

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=909-994

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=3223-3240

TODO: upgrade preserving user edited, dun special case

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2248-2280

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2556-2639

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2930-3032

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=3034-3131

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2590-2593

TODO: com.google.android.carrier service, graphene CarrierConfig2

https://github.com/GrapheneOS/platform_packages_apps_CarrierConfig2/blob/15/src/app/grapheneos/carrierconfig2/loader/Apns.java

TODO: carrier uicc app privs, preserve

https://source.android.com/docs/core/connect/uicc

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2248-2280

## APN Reset

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=5445-5483

## APN Editing

TODO: duplicate check logic

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=5055-5066

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=3132-3186

## CarrierConfig Load

TODO: load path order, ota receiver, telephony content provider, carrierconfig service

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:frameworks/base/services/core/java/com/android/server/updates/CarrierIdInstallReceiver.java

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/services/Telephony/AndroidManifest.xml;l=584-591

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/modules/Permission/tests/cts/permissionpolicy/res/raw/android_manifest.xml;l=4605-4611

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/services/Telephony/src/com/android/phone/CarrierConfigLoader.java;l=227-244

TODO: com.google.android.carrier service, graphene CarrierConfig2

https://github.com/GrapheneOS/platform_packages_apps_CarrierConfig2/blob/15/src/app/grapheneos/carrierconfig2/loader/CarrierConfigLoader.java

TODO: adb override

## com.google.android.carrier protobufs

TODO

### converting to xml

TODO: lossy

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:frameworks/opt/telephony/src/java/com/android/internal/telephony/CarrierResolver.java;l=1054-1106

## Two-digit MNCs

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=5585-5606

https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=299-331

TODO

## Radio HAL

TODO: what's passed to this

https://android.googlesource.com/platform/hardware/interfaces/+/refs/heads/main/radio/1.6/types.hal

## Debugging

TODO: debugging apns using content cmd (list all, list matching, list preferred)

TODO: debugging config using dumpsys carrier_config

TODO: dumping state with dumpsys telephony.registry

TODO: triggering apn reset over adb
