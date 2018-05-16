Syntax:

curl http://127.0.0.1:8081/api/price?token=all
args:
   token is skycoin, bitcoin, samos, all

Example:

```
curl http://127.0.0.1:8081/api/price?token=bitcoin
{
    "ok": 1,
    "data": {
        "bitcoin": {
            "name": "Bitcoin",
            "price_usd": "8417.08",
            "price_btc": "1.0",
            "price_cny": "53599.123732"
        }
    }
}
```
```
curl http://127.0.0.1:8081/aip/price?token=all
{
    "ok": 1,
    "data": {
        "bitcoin": {
            "name": "Bitcoin",
            "price_usd": "8417.08",
            "price_btc": "1.0",
            "price_cny": "53599.123732"
        },
        "samos": {
            "name": "samos",
            "price_usd": "0.1768",
            "price_btc": "0.000021",
            "price_cny": "1.1256"
        },
        "skycoin": {
            "name": "Skycoin",
            "price_usd": "25.0027",
            "price_btc": "0.00297793",
            "price_cny": "159.21469333"
        }
    }
}
```
