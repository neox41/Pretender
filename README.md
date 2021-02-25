# Pretender

Simple reverse proxy to serve web content from a legitimate website, useful for domain categorisation.

```
$>pretender -tls -certificate fullchain.pem -key privkey.pem
2021/02/25 17:47:30 Listening on 443 (TLS)
2021/02/25 17:47:30 Listening on 80
```

Point your domain to the legitimate website

```
Pretender> add google.local https://google.com
2021/02/25 17:49:15 google.local pointing to https://google.com added!
Pretender> add bbc.local https://www.bbc.com
2021/02/25 17:51:40 bbc.local pointing to https://www.bbc.com added!
```


When the client connects to your domain, the legitimate website content will be served

```
Pretender> 2021/02/25 17:52:15 Serving content from https://www.bbc.com
```

