# react-go-spa

POC Single Page App to connect UI (React) / Backend (GO) with any Oauth2 provider `GITHUB`

Optionally you can replace Oauth2 provider very easily in the code by supplying a valid `driver`.

## Server Side Setup

Copy `.env.sample` to `.env`. Replace with valid client id and secret.

> Do not commit changes

## Build 

Run

```bash
make
```

## Heroku

### Buildpack

https://github.com/TV4/go-binary-buildpack.git

### Deployment Process