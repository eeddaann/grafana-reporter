# Grafana Reporter

*Generate reports from grafana automatically*

## configuration
- The ```conf.json``` file contains the configuration.
- The file structered as the example below:
```json
[
{
    "link": "http://127.0.0.1:3000/d/be129f5c-f2ca-4570-8b29-386e61375d28/foo?orgId=1",
    "zoom": "0.25",
    "name": "foo"
},
{
    "link": "http://127.0.0.1:3000/d/ca108980-22b4-4845-aff1-2a2a88380f16/bar?orgId=1",
    "zoom": "0.5",
    "name": "bar"
}
]
```
- **link** - The url to take screenshot from
- **zoom** - The desired zoom (as apears on chrome)
- **name** - The name of the dashboard  