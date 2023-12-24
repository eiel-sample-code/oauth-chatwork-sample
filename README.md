# oauth-chatwork-sample

ChatwokrのOauthをGoで使う場合のサンプルコードです。
一部サンプルレベルの実装になっているので注意してください。

# 利用方法

事前にsetupに書かれた準備をしておいてください。
以下のコマンドでサーバを起動できます。

```
go run main.go
```

サーバーを起動したら https://localhost:8080/ へアクセスしてください。
アクセスするとChatworkのOAuthと接続するための画面が表示されます。
承認すると、localhostへリダイレクトして GET /me を実行しその結果を返します。

# setup

## 必要なツール

* go

devenvを利用している場合は`devenv shell`を実行すれば必要な環境を整えることができます。
また、direnv連携をしたい場合はプロジェクトルートで `devenv init`を実行すれば`.envrc`が生成されます。

## OAuthアプリケーションの用意

Chatworkのサービス連携の画面でOAuthアプリケーションを用意する必要があります。

| クライアント名     | 任意                            |
| クライアントタイプ | confidential                    |
| リダイレクト先URI  | https://localhost:8080/callback |
| スコープ           | users.profile.me:read           |


## 設定する必要がある環境変数

| CHATWORK_OAUTH2_CLIENT_ID | ChatworkのOAuthアプリケーションのクライアントID |
| CHATWORK_OAUTH2_CLIENT_KEY | Chatworkの |

nixを利用している場合は `devenv.local.nix`を作成して定義することで設定できます。
`XXXXX`の部分を置き換えてください。

```nix
{ pkgs, ... }:

{
  env.CHATWORK_OAUTH2_CLIENT_ID = "XXXXXXX";
  env.CHATWORK_OAUTH2_CLIENT_SECRET = "XXXXXX";
}

```

direnvを利用している場合はそこで指定しても構いません

export CHATWORK_OAUTH2_CLIENT_ID="XXXXXXX";
export CHATWORK_OAUTH2_CLIENT_SECRET = "XXXXXX";

## 証明書の作成

チャットワークのOAuthアプリケーションのコールバックURLはHTTPSを指定する必要があるため自己証明を作る必要があります。
mkcertを使う場合は以下で作成できます。
プロジェクトのルートディレクトリで以下のコマンドで作成できます。

```
mkcert -install
mkcert -key-file key.pem -cert-file cert.pem localhost localhost 127.0.0.1 ::1
```

DevEnvを利用している場合はプロジェクトディレクトリで以下のコマンドで作成できます。

```
setup
createCertificate
```
