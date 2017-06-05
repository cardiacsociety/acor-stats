# acor-stats

**Objective**

Generate statistical reports for ACOR that show devices and procedures over time and region.


**Method**

* Import data from spreadsheets / CSV into MongoDB collection
* Provide simple API to extract data for reports
* HTML report page using js charting library

**Doc Design**

There are official registries - devices and procedures, however the devices registry is a subset of procedures so the data can be imported into a single collection and reports generated across both registries.


For a device registry item procType = "device"

```json
{
   "patientId": "2161c1290",
   "siteId": "2161",
   "siteState": "NSW",
   "procDate": "2016-10-01",
   "procType": "device",
   "deviceType": "ICD",
   "deviceSubType": "Dual"
}
```

For a procedures registry item procType = "pci"

```json
{
   "patientId": "2161c1290",
   "siteId": "2161",
   "siteState": "NSW",
   "procDate": "2016-10-01",
   "procType": "pci",
   "deviceType": "stent",
   "deviceSubType": "DES"
}
```

**Demo **

https://acor-report.herokuapp.com/
