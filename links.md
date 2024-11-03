# Links

Links to specs and posts about things useful for when I was troubleshooting this.

TODO: add more stuff to this list from my history

https://www.gsma.com/get-involved/working-groups/wp-content/uploads/2021/01/VoLTE-Implementation-Guide-Jan-2021.pdf

- roaming architectures
- stuff about the well-known "IMS" APN name (i.e., why not old carrier-specific IMS APN names)
- e911 stuff

https://nickvsnetworking.com/best-practices-for-sgw-pgw-deployment-architectures-for-roaming/

- more about roaming

https://www.gsma.com/newsroom/wp-content/uploads/IR.65-v34.0-4.pdf

- more about volte roaming

https://source.android.com/docs/core/connect/ims-service-entitlement

- stuff about how ims provisioning works (and it's not needed for all carriers -- most (all?) canadian ones don't need it)

https://source.android.com/docs/core/connect/ims-single-registration

https://www.gsma.com/solutions-and-impact/technologies/esim/wp-content/uploads/2022/11/SGP.01-v4.3.pdf

- esim provisioning

https://www.etsi.org/deliver/etsi_ts/123100_123199/123122/16.07.00_60/ts_123122v160700p.pdf

- a lot of extremely informative stuff about how network selection works
- describes intended state machines for implementations
- makes it MUCH more clear how PLMNS/EHPLMNS interact, and how it all works when roaming

https://www.gsma.com/newsroom/wp-content/uploads//IR.88-v21.0.pdf

- more about roaming on lte

https://www.etsi.org/deliver/etsi_ts/131100_131199/131102/15.01.00_60/ts_131102v150100p.pdf

- what's stored on the SIM card

https://bluesecblog.wordpress.com/2016/11/24/usim-data-downald-ota-updates/

- remote sim update structure

https://www.etsi.org/deliver/etsi_ts/151000_151099/151011/04.14.00_60/ts_151011v041400p.pdf

- more sim stuff

https://nickvsnetworking.com/sms-transport-supremacy-ip-or-nas/

- sms transport

https://nickvsnetworking.com/bsf-addresses/

- bsf

https://realtimecommunication.wordpress.com/2015/05/27/ut-interface-what-is-it-for/

- xcap/ut

https://www.gsma.com/newsroom/wp-content/uploads/NG.115-v1.0-7.pdf

- iwlan volte

https://nickvsnetworking.com/all-about-ims-authentication-akav1-md5-in-volte-networks/

- aka auth (can be used for ims, passpoint)

https://telecompedia.net/lte-release-causes/

- lte release codes

https://mobilepacketcore.com/eps-bearer-lte/

- eps bearer setup
- makes it more clear where problems are when seeing lte release codes

https://www.3gpp.org/technologies/101-carrier-aggregation-explained

- lte_ca

https://www.3gpp.org/ftp/information/presentations/presentations_2010/2010_06_Latin_America/Hannu_Core%20Network.pdf

- more about lte epc

https://github.com/kyujin-cho/pixel-volte-patch

- Pixel IMS app - gui for CarrierConfig
