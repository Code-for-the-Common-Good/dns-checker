# 🌐dns-checker
## route-example

{API-URL} + /{type}/{dns}/{domain}

##### DNS - options 
- cloudflare
- google
- yandex
- opendns

##### DNS Record types
- a
- aaaa
- cname
- mx
- ns
- ptr
- txt


curl -X GET "{api}/a/cloudflare/google.co.za"
```json
{
    "ipv4": [
        "142.251.47.67"
    ]
}
```
curl -X GET "{api}/aaaa/cloudflare/google.co.za"
```json
{
    "ipv6": [
        "2c0f:fb50:4002:802::2003"
    ]
}
```
curl -X GET "{api}/cname/cloudflare/google.co.za"
```json
{
    "cname": "google.co.za."
}
```
curl -X GET "{api}/mx/cloudflare/google.co.za"
```json
{
    "mx": [
        {
            "Host": "smtp.google.com.",
            "Pref": 0
        }
    ]
}
```
curl -X GET "{api}/ns/cloudflare/google.co.za"
```json
{
    "ns": [
        {
            "Host": "ns4.google.com."
        },
        {
            "Host": "ns1.google.com."
        },
        {
            "Host": "ns3.google.com."
        },
        {
            "Host": "ns2.google.com."
        }
    ]
}
```
curl -X GET "{api}/ptr/cloudflare/"
```json
{
    "ptr": [
        "jnb03s07-in-f3.1e100.net."
    ]
}
```
curl -X GET "{api}/txt/cloudflare/google.co.za"
```json
{
    "txt": [
        "v=spf1 -all"
    ]
}
```