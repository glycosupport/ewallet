
# EWallet

Test assignment for infotecs


## Run Locally

Clone the project

```bash
  git clone https://github.com/glycosupport/ewallet
```

Go to the project directory

```bash
  cd ewallet
```

Build docker compose

```bash
  ./dc-build.sh
```

Up docker compose

```bash
  ./dc-up.sh
```


## API Reference

#### Create wallet

```http
  POST /api/v1/wallet
```

#### Send money

```http
  POST /api/v1/wallet/{walletId}/send
```

#### Get history wallet

```http
  GET /api/v1/wallet/{walletId}/history
```

#### Get state wallet
 
```http
  GET /api/v1/wallet/{walletId}
```