# acor-stats

**Objective**

Generate statistical reports for ACOR that show devices and procedures over time and region.


**Method**

* Import data from spreadsheets / CSV into doc database
* Provide simple API to extract data for reports
* HTML report page using js charting library

**Doc Design**

There are two registries - devices and procedures, however the data is similar for both so it might be possible and practical to model them in a similar way in order to generate statistics across both registries.


For a device registry item can designate proceType as "device"

```json
{
   "patiendId": "2161c1290",
   "siteId": "2161",
   "siteState": "NSW",
   "procDate": "2016-10-01",
   "procType": "device",
   "deviceType": "ICD",
   "deviceSubType": "Dual",
}
```

For a procedures registry item, eg PCI

```json
   "patiendId": "2161c1290",
   "siteId": "2161",
   "siteState": "NSW",
   "procDate": "2016-10-01",
   "procType": "pci",
   "deviceType": "stent",
   "deviceSubType": "DES",
```

