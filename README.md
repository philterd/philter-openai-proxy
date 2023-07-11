# Philter OpenAI Proxy

## Introduction

This project is a proxy for OpenAI that uses [Philter](https://philterd.ai/philter/) to remove PII, PHI, and other sensitive information from a [chat completion](https://platform.openai.com/docs/api-reference/chat) request before sending the request to OpenAI. If you don't have a running instance of Philter, you can launch one in your cloud at https://philterd.ai/philter/.

The proxy works by sending requests destined for OpenAI first to Philter where the sensitive information is redacted per Philter's configuration. The redacted text is then sent to OpenAI. For example, if you send the following text "How old is John Smith?", the proxy and Philter will remove the text "John Smith" from the request. The redacted request sent to OpenAI will be "How old is {{{REDACTED-person}}}?"

Check out the [blog post](https://blog.philterd.ai/removing-pii-phi-from-openai-chat-gpt-api-requests-551f57cef64d) for more information.

## Running the Proxy

```
export PHILTER_ENDPOINT=https://your-philter-ip:8080
./philter-openai-proxy
```

## Usage

To use this proxy, you can send a request to it like you would to OpenAI but change the hostname:

```
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Whose social security number is 123-45-6789"}]
  }'
```

The proxy listens over TLS and requires a certificate and private key. You can generate a self-signed certificate with the following command:

```
make cert
```

## License

Copyright 2023 Philterd, LLC
Licensed under the Apache License, version 2.
