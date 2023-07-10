# Philter OpenAI Proxy

## Introduction

This project is a proxy for OpenAI that uses Philter to remove PII, PHI, and other sensitive information from the request before sending the request to OpenAI.

## Usage

To use this proxy, you can send a request to it like you would to OpenAI but change the hostname:

```
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Who is president of the United States?"}]
  }'
```

The proxy listens over TLS and requires a certificate and private key. You can generate a self-signed certificate with the following command:

```
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname"
```

## License

Copyright 2023 Philterd, LLC
Licensed under the Apache License, version 2.